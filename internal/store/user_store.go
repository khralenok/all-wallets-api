package store

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/all-wallets-api/internal/database"
	"github.com/khralenok/all-wallets-api/internal/models"
)

// Create new user
func AddNewUser(username, passwordHash, baseCurrency string, context *gin.Context) error {
	query := "INSERT INTO users (username, user_pwd, base_currency) VALUES ($1, $2, $3)"

	_, err := database.DB.Exec(query, username, passwordHash, baseCurrency)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to insert user"})
		return err
	}

	return nil
}

// Return user object or error if there is no user with such id in database
func GetUserById(userID int, context *gin.Context) (models.User, error) {
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

// Return user sruct or error if there is no user with such username in database
func GetUserByUsername(username string, context *gin.Context) (models.User, error) {
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

// Return userID(int) of user with provided username. In case if there is no such user, will return -1.
func GetIdByUsername(username string, context *gin.Context) int {
	var userToAdd models.User
	var deleted_at sql.Null[string]

	query := "SELECT * FROM users WHERE username=$1 AND is_deleted = FALSE"

	err := database.DB.QueryRow(query, username).Scan(&userToAdd.ID, &userToAdd.Username, &userToAdd.Password, &userToAdd.BaseCurrency, &userToAdd.CreatedAt, &userToAdd.IsDeleted, &deleted_at)

	if err == sql.ErrNoRows {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "There is no such user"})
		return -1
	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't fetch user data"})
		return -1
	}

	return userToAdd.ID
}

// Return nil if there is no users with such username in database.
func CheckIfUsernameUnique(username string, context *gin.Context) error {
	query := "SELECT * FROM users WHERE username=$1"

	rows, err := database.DB.Query(query, username)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't get response from database"})
		return err
	}

	defer rows.Close()

	if rows.Next() {
		context.JSON(http.StatusConflict, gin.H{"error": "Status Conflict", "message": "This username already taken"})
		return errors.New("status conflict")
	}

	return nil
}

func MarkUserAsDeleted(userID int, context *gin.Context) error {
	query := "UPDATE users SET is_deleted = TRUE, deleted_at = CURRENT_TIMESTAMP WHERE id=$1"

	_, err := database.DB.Exec(query, userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Can't update user data"})
		return err
	}

	return nil
}
