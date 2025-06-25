package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
	"github.com/khralenok/all-wallets-api/store"
)

// Create new wallet with provided name and currency and automatically create new wallet user with admin role based on user who call the function.
func CreateWallet(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var request models.NewWalletRequest
	var newWallet models.Wallet
	var newWalletUser models.WalletUser

	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	query := "INSERT INTO wallets (wallet_name, currency) VALUES ($1, $2) RETURNING id, balance, last_snapshot, created_at"

	err := database.DB.QueryRow(query, request.WalletName, request.Currency).Scan(&newWallet.ID, &newWallet.Balance, &newWallet.LastSnapshot, &newWallet.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert new wallet data into the database"})
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
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert new wallet user data into the database"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_wallet": newWallet, "new_wallet_user": newWalletUser})
}

// Delete the wallet and all it's users. Can be performed only by wallet user with admin role
func DeleteWallet(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	walletID, err := strconv.Atoi(context.Param("wallet_id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !store.CheckUserPermissions(userID, walletID, context) {
		return
	}

	query := "DELETE FROM wallet_users WHERE wallet_id=$1"

	_, err = database.DB.Exec(query, walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Can't delete this wallet users"})
		return
	}

	query = "DELETE FROM wallets WHERE id=$1"

	_, err = database.DB.Exec(query, walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Can't delete this wallet"})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{"status": "No content", "message": "Wallet User was successfully deleted"})
}
