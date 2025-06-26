package logic

import (
	"github.com/khralenok/all-wallets-api/internal/models"
)

// Return sum of all transactions in provided list
func CalcSnapshot(latestTransactions []models.Transaction) int {
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
