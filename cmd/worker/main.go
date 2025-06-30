package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/store"
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
		snapshotCmd.Parse(os.Args[2:])
		latestTransactions, err := store.GetLatestTransactions(*snapshotWalletID)

		if err != nil {
			fmt.Println(err, latestTransactions)
			os.Exit(1)
		}

		sumOfLatestTransactions := logic.CalcSumOfTransactions(latestTransactions)

		err = store.UpdateBalance(*snapshotWalletID, sumOfLatestTransactions)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("balance was successfuly updated")
		os.Exit(0)

	case "xrates":
		rates, err := logic.FetchExchangeRates()

		if err != nil {
			fmt.Println("Error: ", err.Error())
			os.Exit(1)
		}

		for key, value := range rates {
			fmt.Printf("%s: %.2f\n", key, value)
		}

		availableCurrencies, err := store.GetAvailableCurrencies()

		if err != nil {
			fmt.Println("Error: ", err.Error())
			os.Exit(1)
		}

		for _, value := range availableCurrencies {
			fmt.Println(value.Code)
		}

		os.Exit(0)

	default:
		fmt.Println("expected some command")
		os.Exit(1)
	}
}
