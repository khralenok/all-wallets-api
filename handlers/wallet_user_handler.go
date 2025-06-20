package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
)

func ShareWallet(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var request models.NewWalletUserRequest

	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !checkUserPermissions(userID, request.WalletID, context) {
		return
	}

	newUserId := getIdByUsername(request.Username, context)

	if newUserId == -1 {
		return
	}

	if !checkWalletUserUnique(newUserId, request.WalletID, context) {
		return
	}

	if !validateRoleInput(request.UserRole, context) {
		return
	}

	newWalletUser := createWalletUser(request.WalletID, newUserId, request.UserRole, context)

	context.JSON(http.StatusCreated, gin.H{"new_wallet_user": newWalletUser})
}

func checkUserPermissions(userID, walletID int, context *gin.Context) bool {
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

func getIdByUsername(username string, context *gin.Context) int {
	var userToAdd models.User

	query := "SELECT * FROM users WHERE username=$1"

	err := database.DB.QueryRow(query, username).Scan(&userToAdd.ID, &userToAdd.Username, &userToAdd.Password, &userToAdd.BaseCurrency, &userToAdd.CreatedAt)

	if err == sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such user"})
		return -1
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't fetch user data"})
		return -1
	}

	return userToAdd.ID
}

func checkWalletUserUnique(userID, walletID int, context *gin.Context) bool {
	var existingUser models.WalletUser
	query := "SELECT * FROM wallet_users WHERE user_id=$1 and wallet_id=$2"

	err := database.DB.QueryRow(query, userID, walletID).Scan(&existingUser.WalletID, &existingUser.UserID, &existingUser.UserRole, &existingUser.CreatedAt)

	if err != sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "This user already have access to this wallet"})
		return false
	}

	return true
}

func validateRoleInput(userRole string, context *gin.Context) bool {
	validRoles := map[string]bool{"admin": true, "user": true, "spectator": true}

	if !validRoles[userRole] {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "User role value is unacceptable"})
		return false
	}

	return true
}

func createWalletUser(walletID, userID int, userRole string, context *gin.Context) models.WalletUser {
	var newWalletUser models.WalletUser
	query := "INSERT INTO wallet_users (wallet_id, user_id, user_role) VALUES ($1, $2, $3) RETURNING wallet_id, user_id, user_role, created_at"
	err := database.DB.QueryRow(query, walletID, userID, userRole).Scan(&newWalletUser.WalletID, &newWalletUser.UserID, &newWalletUser.UserRole, &newWalletUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert new wallet user data into the database"})
		return models.WalletUser{}
	}

	return newWalletUser
}
