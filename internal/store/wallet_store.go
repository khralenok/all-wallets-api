package store

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

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
