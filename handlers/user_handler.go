package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
	"github.com/khralenok/all-wallets-api/utilities"
)

// Add new user in database or give http error in response. As input require json with username, password and base currency in JSON format
func CreateUser(context *gin.Context) {
	var input models.SigninInputs

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
	}

	if err := checkIfUsernameUnique(input.Username, context); err != nil {
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

	user, err := getUserByUsername(input.Username, context)

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

	user, err := getUserById(userID, context)

	if err != nil {
		return
	}

	var userOutput models.UserOutput

	userOutput.ID = user.ID
	userOutput.Username = user.Username
	userOutput.BaseCurrency = user.BaseCurrency

	userWallets, err := getUserWallets(userID, context)

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

func getUserByUsername(username string, context *gin.Context) (models.User, error) {
	var user models.User
	var deleted_at sql.Null[string]

	query := "SELECT * FROM users WHERE username=$1"
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.BaseCurrency, &user.CreatedAt, &user.IsDeleted, &deleted_at)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return models.User{}, err
	}

	if user.IsDeleted {
		context.JSON(http.StatusGone, gin.H{"error": "Gone", "message": "This user has been deleted"})
		return models.User{}, errors.New("Gone")
	}

	return user, nil
}

func getUserById(userID int, context *gin.Context) (models.User, error) {
	var user models.User
	var deleted_at sql.Null[string]

	query := "SELECT * FROM users WHERE id=$1 AND is_deleted = FALSE"
	err := database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.BaseCurrency, &user.CreatedAt, &user.IsDeleted, &deleted_at)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get user data from database"})
		return models.User{}, err
	}

	return user, nil
}

func getUserWallets(userID int, context *gin.Context) ([]models.WalletOutputSimple, error) {
	var userWallets []models.WalletOutputSimple

	query := "SELECT w.id AS wallet_id, w.wallet_name, w.currency, w.balance, wu.user_role FROM wallets w JOIN wallet_users wu ON wu.wallet_id = w.id WHERE wu.user_id = $1"

	rows, err := database.DB.Query(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get wallets list from database"})
		return userWallets, err
	}

	for rows.Next() {
		var nextWallet models.WalletOutputSimple
		err := rows.Scan(&nextWallet.WalletID, &nextWallet.WalletName, &nextWallet.Currency, &nextWallet.Balance, &nextWallet.UserRole)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't read next row in the wallets list"})
			return userWallets, err
		}

		userWallets = append(userWallets, nextWallet)
	}

	return userWallets, nil
}

func checkIfUsernameUnique(username string, context *gin.Context) error {
	query := "SELECT * FROM users WHERE username=$1"

	rows, err := database.DB.Query(query, username)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get response from database"})
		return err
	}

	if rows.Next() {
		context.JSON(http.StatusConflict, gin.H{"error": "Status Conflict", "message": "This username already taken"})
		return errors.New("status conflict")
	}

	return nil
}
