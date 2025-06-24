package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khralenok/all-wallets-api/database"
	"github.com/khralenok/all-wallets-api/handlers"
	"github.com/khralenok/all-wallets-api/utilities"
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
	router.GET("/profile", utilities.AuthMiddleware(), handlers.GetProfile)
	router.PUT("/delete-user", utilities.AuthMiddleware(), handlers.DeleteUser)

	//Wallets Management
	router.POST("/new-wallet", utilities.AuthMiddleware(), handlers.CreateWallet)
	router.POST("/share-wallet", utilities.AuthMiddleware(), handlers.AddWalletUser)
	router.DELETE("/remove-wallet-user/wallet/:wallet_id/username/:username/", utilities.AuthMiddleware(), handlers.DeleteWalletUser)
	router.DELETE("/delete-wallet/:wallet_id", utilities.AuthMiddleware(), handlers.DeleteWallet)

	router.Run(":8080")
}
