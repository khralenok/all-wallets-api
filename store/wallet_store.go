package store

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
)

// Return true if there is wallet with such id in database
func IsWalletExist(id int, context *gin.Context) bool {
	var wallet models.Wallet

	query := "SELECT 1 FROM wallets WHERE id = $1"

	err := database.DB.QueryRow(query, id).Scan(&wallet.ID, &wallet.WalletName, &wallet.Currency, &wallet.Balance, &wallet.LastSnapshot, &wallet.CreatedAt)

	if err == sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such wallet"})
		return false
	}

	return true
}
