package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
	"github.com/khralenok/all-wallets-api/store"
)

// Create new wallet user with provided role, based on target wallet id and username of user who need to be added. Can be permormed only by wallet admin.
func AddWalletUser(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var request models.NewWalletUserRequest

	if err := context.BindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return
	}

	if !store.CheckUserPermissions(userID, request.WalletID, context) {
		return
	}

	newUserId := store.GetIdByUsername(request.Username, context)

	if newUserId == -1 {
		return
	}

	if !store.CheckWalletUserUnique(newUserId, request.WalletID, context) {
		return
	}

	if !validateRoleInput(request.UserRole, context) {
		return
	}

	newWalletUser := store.CreateWalletUser(request.WalletID, newUserId, request.UserRole, context)

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

	query := "DELETE FROM wallet_users WHERE user_id=$1 AND wallet_id=$2"

	_, err = database.DB.Exec(query, userToDeleteID, walletID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "raw_error": err.Error(), "message": "Can't delete this user"})
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
