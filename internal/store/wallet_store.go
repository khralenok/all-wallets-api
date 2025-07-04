package store

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/logic"
	"github.com/khralenok/all-wallets-api/internal/models"
)

func AddNewWallet(input models.NewWalletRequest, context *gin.Context) (models.Wallet, error) {
	var newWallet models.Wallet
	query := "INSERT INTO wallets (wallet_name, currency) VALUES ($1, $2) RETURNING *"

	err := database.DB.QueryRow(query, input.WalletName, input.Currency).Scan(&newWallet.ID, &newWallet.WalletName, &newWallet.Currency, &newWallet.Balance, &newWallet.LastSnapshot, &newWallet.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert new wallet data into the database"})
		return models.Wallet{}, err
	}

	return newWallet, nil
}

// Remove specified wallet from database
func RemoveWallet(walletID int, context *gin.Context) error {
	query := "DELETE FROM wallets WHERE id=$1"

	_, err := database.DB.Exec(query, walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't delete this wallet"})
		return err
	}

	return nil
}

// Return true if there is wallet with such id in database
func IsWalletExist(id int, context *gin.Context) bool {
	var walletID int

	query := "SELECT 1 FROM wallets WHERE id = $1"

	err := database.DB.QueryRow(query, id).Scan(&walletID)

	if err == sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such wallet"})
		return false
	}

	return true
}

// Return wallet struct by provided id or an error
func GetWalletByID(walletID int) (models.Wallet, error) {
	var wallet models.Wallet

	query := "SELECT * FROM wallets WHERE id = $1"

	err := database.DB.QueryRow(query, walletID).Scan(&wallet.ID, &wallet.WalletName, &wallet.Currency, &wallet.Balance, &wallet.LastSnapshot, &wallet.CreatedAt)

	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

// Return array of simplified wallet structs or an error.
func GetWalletsByUser(user models.User, context *gin.Context) ([]models.WalletOutput, error) {
	var userWallets []models.WalletOutput

	userCurrencyDecimalPlaces, err := GetCurrencyDecimalPlaces(user.BaseCurrency)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get user currency metadata"})
		return userWallets, err
	}

	query := "SELECT w.id AS wallet_id, w.wallet_name, w.currency, w.balance, wu.user_role, cm.decimal_places FROM wallets w JOIN wallet_users wu ON wu.wallet_id = w.id JOIN currency_metadata cm ON w.currency = cm.code WHERE wu.user_id = $1"

	rows, err := database.DB.Query(query, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get wallets list from database"})
		return userWallets, err
	}

	defer rows.Close()

	for rows.Next() {
		var nextWallet models.WalletOutput
		var rawBalance int
		var decimalPlaces int

		err := rows.Scan(&nextWallet.WalletID, &nextWallet.WalletName, &nextWallet.Currency, &rawBalance, &nextWallet.UserRole, &decimalPlaces)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't read next row in the wallets list"})
			return []models.WalletOutput{}, err
		}

		latestTransactions, err := GetLatestTransactions(nextWallet.WalletID)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get transactions data for this wallet"})
			return []models.WalletOutput{}, err
		}

		rawBalance += logic.CalcSumOfTransactions(latestTransactions)

		nextWallet.Balance = logic.FormatOutputValue(rawBalance, decimalPlaces)

		if nextWallet.Currency != user.BaseCurrency {

			rate, err := GetRate(nextWallet.Currency, user.BaseCurrency, context)

			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't fetch exchange rates"})
				return []models.WalletOutput{}, err
			}

			balanceToExchange, err := strconv.ParseFloat(nextWallet.Balance, 64)

			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Some problem with exchange rates occurred"})
				return []models.WalletOutput{}, err
			}

			nextWallet.UserCurrencyBalance = strconv.FormatFloat(balanceToExchange*rate, 'f', userCurrencyDecimalPlaces, 64)
		} else {
			nextWallet.UserCurrencyBalance = nextWallet.Balance
		}

		userWallets = append(userWallets, nextWallet)
	}

	return userWallets, nil
}

// Return decimal places for specific wallet. Don't try to use it for definining currency specific decimal places — instead use specialized function GetCurrencyDecimalPlaces.
func GetWalletDecimalPlaces(walletID int) (int, error) {
	var walletCurrency string
	var decimalPlaces int

	query := "SELECT cm.code, cm.decimal_places FROM wallets w JOIN currency_metadata cm ON cm.code = w.currency WHERE w.id = $1"

	err := database.DB.QueryRow(query, walletID).Scan(&walletCurrency, &decimalPlaces)

	if err != nil {
		return -1, err
	}

	return decimalPlaces, nil
}

// Update balance for specified sum.
func UpdateBalance(walletID, sumOfLatestTransactions int) error {
	wallet, err := GetWalletByID(walletID)

	if err != nil {
		return err
	}

	newBalance := wallet.Balance + sumOfLatestTransactions
	newSnapshotTime := time.Now()

	query := "UPDATE wallets SET balance = $1, last_snapshot = $2 WHERE id = $3"

	_, err = database.DB.Exec(query, newBalance, newSnapshotTime, walletID)

	if err != nil {
		return err
	}

	return nil
}

// Return true if considering latest transactions wallet have enough funds to add expense of specified amount
func CheckIfBalanceIsEnough(walletID, expense int) (bool, error) {
	wallet, err := GetWalletByID(walletID)

	if err != nil {
		return false, err
	}

	latestTransactions, err := GetLatestTransactions(walletID)

	if err != nil {
		return false, err
	}

	currentBalance := wallet.Balance + logic.CalcSumOfTransactions(latestTransactions)

	if expense > currentBalance {
		return false, nil
	}

	return true, nil
}
