package commands

import (
	"flag"
	"os"

	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Worker function for updating wallet snapshot by sum up all new trunsactions to balance stored in wallet balance
func UpdateWalletSnapshot(snapshotCmd *flag.FlagSet, snapshotWalletID *int) error {
	snapshotCmd.Parse(os.Args[2:])
	latestTransactions, err := store.GetLatestTransactions(*snapshotWalletID)

	if err != nil {
		return err
	}

	sumOfLatestTransactions := logic.CalcSumOfTransactions(latestTransactions)

	err = store.UpdateBalance(*snapshotWalletID, sumOfLatestTransactions)

	if err != nil {
		return err
	}

	return nil
}
