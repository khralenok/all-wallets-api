package store

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

func AddTransaction(input models.TransactionInput, isDeposit bool, context *gin.Context) (models.Transaction, error) {
	userID := context.MustGet("userID").(int)
	var newTransaction models.Transaction

	query := "INSERT INTO transactions (amount, is_deposit, category, wallet_id, creator_id) VALUES ($1, $2, $3, $4, $5) RETURNING *"

	err := database.DB.QueryRow(query, input.Amount, isDeposit, input.Category, input.WalletID, userID).Scan(&newTransaction.ID, &newTransaction.Amount, &newTransaction.IsDeposit, &newTransaction.Category, &newTransaction.WalletID, &newTransaction.CreatorID, &newTransaction.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Failed to insert new transaction data into the database"})
		return models.Transaction{}, err
	}

	return newTransaction, nil
}
