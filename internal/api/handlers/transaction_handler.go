package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Create new transaction
func CreateTransaction(context *gin.Context, isDeposit bool) {
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

	if !isDeposit {
		if isAllowed, err := store.CheckIfBalanceIsEnough(input.WalletID, formatedAmount); !isAllowed {
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't approve wallet have enough funds to write expense"})
				return
			}

			context.JSON(http.StatusConflict, gin.H{"error": "Conflict", "message": "Insufficient funds. You cannot spend more than your current wallet balance."})
			return
		}
	}

	newTransaction, err := store.AddTransaction(formatedAmount, input.WalletID, isDeposit, input.Category, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_transaction": newTransaction})
}

// Get user transactions
func GetWalletTransactions(context *gin.Context) {
	_ = context.MustGet("userID").(int)

	walletID, err := strconv.Atoi(context.Param("wallet_id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	transactionsRaw, err := store.GetWalletTransactions(walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't fetch the data from database"})
		return
	}

	decimalPlaces, err := store.GetWalletDecimalPlaces(walletID)

	if err != nil {
		return
	}

	transactions := logic.FormatTransactionOutput(transactionsRaw, decimalPlaces)

	context.JSON(http.StatusOK, gin.H{"status": "OK", "transactions": transactions})
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
