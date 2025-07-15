package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Create new wallet with provided name and currency and automatically create new wallet user with admin role based on user who call the function.
func CreateWallet(context *gin.Context) {
	userId := context.MustGet("userID").(int)

	var input models.NewWalletRequest

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	newWallet, err := store.AddNewWallet(input, context)

	if err != nil {
		return
	}

	newWalletUser, err := store.AddWalletUser(newWallet.ID, userId, "admin", context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_wallet": newWallet, "new_wallet_user": newWalletUser})
}

// Response with wallet data
func GetWallet(context *gin.Context) {
	_ = context.MustGet("userID").(int)
	walletID, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Id parameter should be integer"})
		return
	}

	// TO DO: Check if user have rights to see the wallet

	wallet, err := store.GetWalletByID(walletID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Can't find such wallet"})
		return
	}

	latestTransactions, err := store.GetLatestTransactions(walletID)

	if err != nil {
		return
	}

	decimalPlaces, err := store.GetWalletDecimalPlaces(walletID)

	if err != nil {
		return
	}

	balance := wallet.Balance + logic.CalcSumOfTransactions(latestTransactions)

	outputBalance := logic.FormatOutputValue(balance, decimalPlaces)

	context.JSON(http.StatusOK, gin.H{"id": wallet.ID, "wallet_name": wallet.WalletName, "balance": outputBalance, "currency": wallet.Currency, "last_snapshot": wallet.LastSnapshot, "created_at": wallet.CreatedAt})
}

// Delete the wallet and all it's users. Can be performed only by wallet user with admin role
func DeleteWallet(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	walletID, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !store.CheckUserPermissions(userID, walletID, context) {
		return
	}

	if err := store.RemoveWalletUser(walletID, userID, context); err != nil {
		return
	}

	if err := store.RemoveWallet(walletID, context); err != nil {
		return
	}

	context.JSON(http.StatusNoContent, gin.H{"status": "No content", "message": "Wallet User was successfully deleted"})
}
