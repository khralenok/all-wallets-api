package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/api/middleware"
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

	passwordHash, err := middleware.HashPassword(strings.TrimSpace(input.Password))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Password encryption failed"})
		return
	}

	err = store.AddNewUser(input.Username, passwordHash, input.BaseCurrency, context)

	if err != nil {
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

	if !middleware.CheckPasswordHash(strings.TrimSpace(input.Password), strings.TrimSpace(user.Password)) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Status Unauthorized", "message": "Invalid credentials"})
		return
	}

	token, err := middleware.GenerateJWT(user.ID)

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

// Mark the user as deleted and remove them from all wallets they participated in. Can be called only by user themself.
func DeleteUser(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	if store.DeleteUserFromAllWallets(userID, context) != nil {
		return
	}

	if store.MarkUserAsDeleted(userID, context) != nil {
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "User was successfully deleted"})
}
