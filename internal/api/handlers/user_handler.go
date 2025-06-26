package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	utilities "github.com/khralenok/all-wallets-api/internal/api/middleware"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Add new user in database or give http error in response. As input require json with username, password and base currency in JSON format
func CreateUser(context *gin.Context) {
	var input models.SigninInputs

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
	}

	if err := store.CheckIfUsernameUnique(input.Username, context); err != nil {
		return
	}

	passwordHash, err := utilities.HashPassword(strings.TrimSpace(input.Password))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Password encryption failed"})
		return
	}

	query := "INSERT INTO users (username, user_pwd, base_currency) VALUES ($1, $2, $3)"

	_, err = database.DB.Exec(query, input.Username, passwordHash, input.BaseCurrency)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"status": "Created"})
}

// Response with JWT Token user need to use other API handlers. As input require username and password in JSON format
func LoginUser(context *gin.Context) {
	var input models.LoginInputs

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	user, err := store.GetUserByUsername(input.Username, context)

	if err != nil {
		return
	}

	if !utilities.CheckPasswordHash(strings.TrimSpace(input.Password), strings.TrimSpace(user.Password)) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Status Unauthorized", "message": "Invalid credentials"})
		return
	}

	token, err := utilities.GenerateJWT(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Token generation failed"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "token": token})
}

// Response with user data and brief data for all wallets which user participated in.
func GetProfile(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	user, err := store.GetUserById(userID, context)

	if err != nil {
		return
	}

	var userOutput models.UserOutput

	userOutput.ID = user.ID
	userOutput.Username = user.Username
	userOutput.BaseCurrency = user.BaseCurrency

	userWallets, err := store.GetUserWallets(userID, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": userOutput, "wallets": userWallets})
}

// Mark the user as deleted and remove them from all wallets they participated in. Can be performed only by user themself.
func DeleteUser(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	query := "DELETE FROM wallet_users WHERE user_id=$1"

	_, err := database.DB.Exec(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't remove this user from wallet users list"})
		return
	}

	query = "UPDATE users SET is_deleted = TRUE, deleted_at = CURRENT_TIMESTAMP WHERE id=$1"

	_, err = database.DB.Exec(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't update user data"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "User was successfully deleted"})
}
