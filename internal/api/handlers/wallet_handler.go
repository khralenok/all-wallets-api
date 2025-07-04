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
	_ = context.MustGet("userID").(int)

	var input models.NewWalletRequest
	var newWalletUserRequest models.NewWalletUserRequest

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	newWallet, err := store.AddNewWallet(input, context)

	if err != nil {
		return
	}

	newWalletUserRequest.WalletID = newWallet.ID
	newWalletUserRequest.UserRole = "admin"

	newWalletUser, err := store.AddWalletUser(newWalletUserRequest, context)

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

	// TO DO: Replace with correct output format
	context.JSON(http.StatusOK, gin.H{"id": wallet.ID, "balance": outputBalance, "currency": wallet.Currency})
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
