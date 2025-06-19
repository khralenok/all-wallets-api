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

	var supplicant models.WalletUser
	var newWalletUser models.WalletUser
	var request models.NewWalletUserRequest

	//1. Get inputs

	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	//4. Check if user who request adding have enough permissions and wallet exists

	query := "SELECT * FROM wallet_users WHERE user_id=$1 and wallet_id=$2"

	err := database.DB.QueryRow(query, userID, request.WalletID).Scan(&supplicant.WalletID, &supplicant.UserID, &supplicant.UserRole, &supplicant.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to fetch supplicant user data"})
		return
	}

	if supplicant.UserRole != "admin" {
		context.JSON(http.StatusForbidden, gin.H{"error": "Status Forbidden", "message": "You have not enough no permission to add new users"})
		return
	}

	//3. Check if user we should add exist and not in the list of users already

	var userToAdd models.User

	query = "SELECT * FROM users WHERE username=$1"

	err = database.DB.QueryRow(query, request.Username).Scan(&userToAdd.ID, &userToAdd.Username, &userToAdd.Password, &userToAdd.BaseCurrency, &userToAdd.CreatedAt)

	if err == sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such user"})
		return
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't fetch user data"})
		return
	}

	var existingUser models.WalletUser

	query = "SELECT * FROM wallet_users WHERE user_id=$1 and wallet_id=$2"

	err = database.DB.QueryRow(query, userToAdd.ID, request.WalletID).Scan(&existingUser.WalletID, &existingUser.UserID, &existingUser.UserRole, &existingUser.CreatedAt)

	if err != sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "This user already have access to this wallet"})
		return
	}

	//4. Check if role value is appropriate

	validRoles := map[string]bool{"admin": true, "user": true, "spectator": true}

	if !validRoles[request.UserRole] {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "User role value is unacceptable"})
		return
	}

	//5. Create new row in wallet users with corresponding data
	query = "INSERT INTO wallet_users (wallet_id, user_id, user_role) VALUES ($1, $2, $3) RETURNING wallet_id, user_id, user_role, created_at"
	err = database.DB.QueryRow(query, request.WalletID, userToAdd.ID, request.UserRole).Scan(&newWalletUser.WalletID, &newWalletUser.UserID, &newWalletUser.UserRole, &newWalletUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "detail": userToAdd.ID, "message": "Failed to insert new wallet user data into the database"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_wallet_user": newWalletUser})
}
