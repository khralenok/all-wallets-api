package store

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

// Insert new wallet user to database. Return new wallet user sruct.
func AddWalletUser(walletId, userId int, userRole string, context *gin.Context) (models.WalletUser, error) {
	var newWalletUser models.WalletUser

	query := "INSERT INTO wallet_users (wallet_id, user_id, user_role) VALUES ($1, $2, $3) RETURNING *"
	err := database.DB.QueryRow(query, walletId, userId, userRole).Scan(&newWalletUser.WalletID, &newWalletUser.UserID, &newWalletUser.UserRole, &newWalletUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert new wallet user data into the database"})
		return models.WalletUser{}, err
	}

	return newWalletUser, nil
}

// Check if there in no user with such user id in specific wallet users list. Return true if there is no user with such id.
func CheckWalletUserUnique(userID, walletID int, context *gin.Context) bool {
	var existingUser models.WalletUser
	query := "SELECT * FROM wallet_users WHERE user_id=$1 and wallet_id=$2"

	err := database.DB.QueryRow(query, userID, walletID).Scan(&existingUser.WalletID, &existingUser.UserID, &existingUser.UserRole, &existingUser.CreatedAt)

	if err != sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "This user already have access to this wallet"})
		return false
	}

	return true
}

// Return true if user have admin role for wallet with provided ID.
func CheckUserPermissions(userID, walletID int, context *gin.Context) bool {
	var supplicant models.WalletUser

	query := "SELECT * FROM wallet_users WHERE user_id=$1 and wallet_id=$2"

	err := database.DB.QueryRow(query, userID, walletID).Scan(&supplicant.WalletID, &supplicant.UserID, &supplicant.UserRole, &supplicant.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to fetch supplicant user data"})
		return false
	}

	if supplicant.UserRole != "admin" {
		context.JSON(http.StatusForbidden, gin.H{"error": "Status Forbidden", "message": "You have not enough no permission to add new users"})
		return false
	}

	return true
}

// Remove the row with such user_id and wallet_id from wallet_users table
func RemoveWalletUser(walletID, userToDeleteId int, context *gin.Context) error {
	query := "DELETE FROM wallet_users WHERE user_id=$1 AND wallet_id=$2"

	_, err := database.DB.Exec(query, userToDeleteId, walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't delete this user"})
		return err
	}

	context.JSON(http.StatusNoContent, gin.H{"status": "No content", "message": "Wallet User was successfully deleted"})
	return nil
}

// Remove all rows with this user from wallet_users table
func RemoveUserFromAllWallets(userID int, context *gin.Context) error {
	query := "DELETE FROM wallet_users WHERE user_id=$1"

	_, err := database.DB.Exec(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't remove this user from wallet users list"})
		return err
	}

	return nil
}
