package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/models"
	"github.com/khralenok/all-wallets-api/internal/store"
)

// Create new wallet user with provided role, based on target wallet id and username of user who need to be added. Can be permormed only by wallet admin.
func CreateWalletUser(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var input models.NewWalletUserRequest

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !store.CheckUserPermissions(userID, input.WalletID, context) {
		return
	}

	newUserId := store.GetIdByUsername(input.Username, context)

	if newUserId == -1 {
		return
	}

	if !store.CheckWalletUserUnique(newUserId, input.WalletID, context) {
		return
	}

	if !validateRoleInput(input.UserRole, context) {
		return
	}

	newWalletUser, err := store.AddWalletUser(input.WalletID, newUserId, input.UserRole, context)

	if err != nil {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"new_wallet_user": newWalletUser})
}

// Remove user from list of wallet users, so it can gain access to wallet anymore. Wallet ID and Username must be provided via url query
func DeleteWalletUser(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	walletID, err := strconv.Atoi(context.Param("wallet_id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	username := context.Param("username")

	if username == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !store.CheckUserPermissions(userID, walletID, context) {
		return
	}

	userToDeleteID := store.GetIdByUsername(username, context)

	if userToDeleteID == -1 {
		return
	}

	if err := store.RemoveWalletUser(walletID, userToDeleteID, context); err != nil {
		return
	}

	context.JSON(http.StatusNoContent, gin.H{"status": "No content", "message": "Wallet User was successfully deleted"})
}

func validateRoleInput(userRole string, context *gin.Context) bool {
	validRoles := map[string]bool{"admin": true, "user": true, "spectator": true}

	if !validRoles[userRole] {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "User role value is unacceptable"})
		return false
	}

	return true
}
