package logic

import (
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Return sum of all transactions in provided list
func CalcSumOfTransactions(latestTransactions []models.Transaction) int {
	var sum int

	for i := 0; i < len(latestTransactions); i++ {
		multiplier := 1

		if !latestTransactions[i].IsDeposit {
			multiplier = -1
		}

		sum += latestTransactions[i].Amount * multiplier
	}

	return sum
}

// Return true if considering latest transactions wallet have enough funds to add expense of specified amount
func CheckIfBalanceIsEnough(walletID, expense int) (bool, error) {
	wallet, err := store.GetWalletByID(walletID)

	if err != nil {
		return false, err
	}

	latestTransactions, err := store.GetLatestTransactions(walletID)

	if err != nil {
		return false, err
	}

	currentBalance := wallet.Balance + CalcSumOfTransactions(latestTransactions)

	if expense > currentBalance {
		return false, nil
	}

	return true, nil
}
