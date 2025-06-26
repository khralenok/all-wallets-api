package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

func AddIncome(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	input, err := validateTransactionInput(userID, context)

	if err != nil {
		return
	}

	newTransaction, err := store.AddTransaction(input, true, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_transaction": newTransaction})
}

func AddExpense(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	input, err := validateTransactionInput(userID, context)

	if err != nil {
		return
	}

	newTransaction, err := store.AddTransaction(input, false, context)

	if err != nil {
		return
	}

	// TO DO: Check if wallet have enough funds

	context.JSON(http.StatusCreated, gin.H{"new_transaction": newTransaction})
}

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
