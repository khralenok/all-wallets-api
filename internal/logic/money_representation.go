package logic

import (
	"fmt"
	"math"

	"github.com/khralenok/all-wallets-api/internal/models"
)

// Take amount in human readable format and turn it to format convinient for system level operations
func FormatInputValue(amount float64, decimalPlaces int) int {
	multiplier := math.Pow10(decimalPlaces)

	return int(math.Round(amount * multiplier))
}

// Take amount in convinient for system level operations format and turn it to human readable format
func FormatOutputValue(amount int, decimalPlaces int) string {
	divisor := math.Pow10(decimalPlaces)

	return fmt.Sprintf("%.*f", decimalPlaces, float64(float64(amount)/divisor))
}

func FormatTransactionOutput(rawTransactions []models.Transaction, decimalPlaces int) []models.TransactionOutput {
	var transactions []models.TransactionOutput

	for _, rawTrx := range rawTransactions {
		var newTrx models.TransactionOutput

		newTrx.ID = rawTrx.ID
		newTrx.IsDeposit = rawTrx.IsDeposit
		newTrx.Category = rawTrx.Category
		newTrx.CreatorID = rawTrx.CreatorID
		newTrx.CreatedAt = rawTrx.CreatedAt

		newTrx.Amount = FormatOutputValue(rawTrx.Amount, decimalPlaces)

		transactions = append(transactions, newTrx)
	}

	return transactions
}
