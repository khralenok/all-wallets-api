package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Add income transaction of specified amount and category to provided wallet
func AddIncome(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	input, err := validateTransactionInput(userID, context)

	if err != nil {
		return
	}

	decimalPlaces, err := store.GetWalletDecimalPlaces(input.WalletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Failed to get decimal places"})
		return
	}

	formatedAmount := logic.FormatInputValue(input.Amount, decimalPlaces)

	newTransaction, err := store.AddTransaction(formatedAmount, input.WalletID, true, input.Category, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_transaction": newTransaction})
}

// Add expense transaction of specified amount and category to provided wallet
func AddExpense(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	input, err := validateTransactionInput(userID, context)

	if err != nil {
		return
	}

	decimalPlaces, err := store.GetWalletDecimalPlaces(input.WalletID)

	if err != nil {
		return
	}

	formatedAmount := logic.FormatInputValue(input.Amount, decimalPlaces)

	isAllowed, err := logic.CheckIfBalanceIsEnough(input.WalletID, formatedAmount)

	if !isAllowed {
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't approve wallet have enough funds to write expense"})
			return
		}

		context.JSON(http.StatusConflict, gin.H{"error": "Conflict", "message": "Insufficient funds. You cannot spend more than your current wallet balance."})
		return
	}

	newTransaction, err := store.AddTransaction(formatedAmount, input.WalletID, true, input.Category, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_transaction": newTransaction})
}

// Group checkups that are common for adding expense and income based on user input
func validateTransactionInput(userID int, context *gin.Context) (models.TransactionInput, error) {

	var input models.TransactionInput

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return models.TransactionInput{}, err
	}

	if input.Amount <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Amount must be positive"})
		return models.TransactionInput{}, errors.New("amount must be positive")
	}

	if !store.IsWalletExist(input.WalletID, context) {
		return models.TransactionInput{}, errors.New("seems wallet doesn't exist")
	}

	if !store.CheckUserPermissions(userID, input.WalletID, context) {
		return models.TransactionInput{}, errors.New("user have no rights to add transactions")
	}

	return input, nil
}
