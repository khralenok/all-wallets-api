package store

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

// Add new transaction to DB. Return transaction object
func AddTransaction(amount, walletID int, isDeposit bool, category string, context *gin.Context) (models.Transaction, error) {
	userID := context.MustGet("userID").(int)
	var newTransaction models.Transaction

	query := "INSERT INTO transactions (amount, is_deposit, category, wallet_id, creator_id, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"

	err := database.DB.QueryRow(query, amount, isDeposit, category, walletID, userID, time.Now()).Scan(&newTransaction.ID, &newTransaction.Amount, &newTransaction.IsDeposit, &newTransaction.Category, &newTransaction.WalletID, &newTransaction.CreatorID, &newTransaction.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Failed to insert new transaction data into the database"})
		return models.Transaction{}, err
	}

	return newTransaction, nil
}

// Return array of all transactions from lastest snapshot
func GetLatestTransactions(walletID int) ([]models.Transaction, error) {
	var latestTransactions []models.Transaction

	wallet, err := GetWalletByID(walletID)

	if err != nil {
		return []models.Transaction{}, err
	}

	query := "SELECT * FROM transactions WHERE wallet_id = $1 AND created_at > $2"

	rows, err := database.DB.Query(query, wallet.ID, wallet.LastSnapshot)

	if err != nil {
		return []models.Transaction{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var nextTransaction models.Transaction
		err := rows.Scan(&nextTransaction.ID, &nextTransaction.Amount, &nextTransaction.IsDeposit, &nextTransaction.Category, &nextTransaction.WalletID, &nextTransaction.CreatorID, &nextTransaction.CreatedAt)

		if err != nil {
			return []models.Transaction{}, err
		}

		latestTransactions = append(latestTransactions, nextTransaction)
	}

	return latestTransactions, nil
}
