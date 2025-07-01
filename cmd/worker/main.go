package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/khralenok/all-wallets-api/internal/commands"
	"github.com/khralenok/all-wallets-api/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer database.DB.Close()

	snapshotCmd := flag.NewFlagSet("snapshot", flag.ExitOnError)
	snapshotWalletID := snapshotCmd.Int("id", 0, "Id of wallet you want to make snapshot for")

	switch os.Args[1] {
	case "snapshot":
		err := commands.UpdateWalletSnapshot(snapshotCmd, snapshotWalletID)

		if err != nil {
			fmt.Println("Error: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Balance was successfuly updated!")
		os.Exit(0)

	case "xrates":
		err := commands.UpdateExchangeRates()

		if err != nil {
			fmt.Println("Error: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Exchange rates were successfuly updated!")
		os.Exit(0)

	default:
		fmt.Println("expected some command")
		os.Exit(1)
	}
}
