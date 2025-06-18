package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
)

func CreateWallet(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	//1. Validate inputs +
	//2. Create new wallet +
	//3. Create new wallet user
	//4. Return success message

	var request models.NewWalletRequest
	var newWallet models.Wallet
	var newWalletUser models.WalletUser

	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "INSERT INTO wallets (wallet_name, currency) VALUES ($1, $2) RETURNING id, balance, last_snapshot, created_at"

	err := database.DB.QueryRow(query, request.WalletName, request.Currency).Scan(&newWallet.ID, &newWallet.Balance, &newWallet.LastSnapshot, &newWallet.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Failed to insert new wallet data into the database"})
		return
	}

	newWallet.WalletName = request.WalletName
	newWallet.Currency = request.Currency

	newWalletUser.UserID = userID
	newWalletUser.WalletID = newWallet.ID
	newWalletUser.UserRole = "admin"

	query = "INSERT INTO wallet_users (wallet_id, user_id, user_role) VALUES ($1, $2, $3) RETURNING created_at"

	err = database.DB.QueryRow(query, newWalletUser.WalletID, newWalletUser.UserID, newWalletUser.UserRole).Scan(&newWalletUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Failed to insert new wallet user data into the database"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_wallet": newWallet, "new_wallet_user": newWalletUser})
}
