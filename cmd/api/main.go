package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khralenok/all-wallets-api/internal/api/handlers"
	"github.com/khralenok/all-wallets-api/internal/api/middleware"
	"github.com/khralenok/all-wallets-api/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer database.DB.Close()

	router := gin.Default()

	//User Management
	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)
	router.PUT("/delete-user", middleware.AuthMiddleware(), handlers.DeleteUser)

	//Wallets Management
	router.POST("/new-wallet", middleware.AuthMiddleware(), handlers.CreateWallet)
	router.POST("/share-wallet", middleware.AuthMiddleware(), handlers.AddWalletUser)
	router.DELETE("/remove-wallet-user/wallet/:wallet_id/username/:username/", middleware.AuthMiddleware(), handlers.DeleteWalletUser)
	router.DELETE("/delete-wallet/:wallet_id", middleware.AuthMiddleware(), handlers.DeleteWallet)

	//Transactions Management
	router.POST("/add-income", middleware.AuthMiddleware(), handlers.AddIncome)
	router.POST("/add-expense", middleware.AuthMiddleware(), handlers.AddExpense)

	router.Run(":8080")
}
