package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/models"
	"github.com/khralenok/all-wallets-api/utilities"
)

func CreateUser(context *gin.Context) {
	var newUser models.User

	if err := context.BindJSON(&newUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
	}

	query := "SELECT * FROM users WHERE username=$1"

	rows, err := database.DB.Query(query, newUser.Username)

	if err != nil {
		panic(err.Error())
	}

	if rows.Next() {
		context.JSON(http.StatusConflict, gin.H{"error": "Status Conflict", "message": "This username already taken"})
		return
	}

	var passwordHash string

	if passwordHash, err = utilities.HashPassword(newUser.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Password encryption failed"})
		return
	}

	query = "INSERT INTO users (username, user_pwd, base_currency) VALUES ($1, $2, $3) RETURNING id, created_at"

	err = database.DB.QueryRow(query, newUser.Username, passwordHash, newUser.BaseCurrency).Scan(&newUser.ID, &newUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert user"})
		return
	}

	context.JSON(http.StatusCreated, newUser)
}

func LoginUser(context *gin.Context) {
	var input models.LoginInputs

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid input format"})
	}

	var user models.User
	var err error

	if user, err = getUserByUsername(input.Username, context); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such user"})
		return
	}

	if !utilities.CheckPasswordHash(input.Password, user.Password) {
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

func GetProfile(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var user models.User

	query := "SELECT * FROM users WHERE id=$1"
	err := database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.BaseCurrency, &user.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get user data from database"})
		return
	}

	var userWallets []models.Wallet

	query = "SELECT w.* FROM wallets w JOIN wallet_users wu ON wu.wallet_id = w.id WHERE wu.user_id = $1"

	rows, err := database.DB.Query(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get wallets list from database"})
		return
	}

	for rows.Next() {
		var nextWallet models.Wallet
		err := rows.Scan(&nextWallet.ID, &nextWallet.WalletName, &nextWallet.Currency, &nextWallet.Balance, &nextWallet.LastSnapshot, &nextWallet.CreatedAt)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't read next row in the wallets list"})
			return
		}

		userWallets = append(userWallets, nextWallet)
	}

	context.JSON(http.StatusOK, gin.H{"user": user, "wallets": userWallets})
}

func getUserByUsername(username string, context *gin.Context) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username=$1"
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.BaseCurrency, &user.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Bad Request", "message": "Invalid input format"})
		return models.User{}, err
	}

	return user, nil
}
