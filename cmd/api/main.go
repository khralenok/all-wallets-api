package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//User Management
	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)
	router.PUT("/delete-user", middleware.AuthMiddleware(), handlers.DeleteUser)

	//Wallets Management
	router.POST("/new-wallet", middleware.AuthMiddleware(), handlers.CreateWallet)
	router.GET("/wallet/:id", middleware.AuthMiddleware(), handlers.GetWallet)
	router.DELETE("/delete-wallet/:id", middleware.AuthMiddleware(), handlers.DeleteWallet)

	//Wallet Users Management
	router.POST("/share-wallet", middleware.AuthMiddleware(), handlers.CreateWalletUser)
	router.DELETE("/remove-wallet-user/wallet/:wallet_id/username/:username/", middleware.AuthMiddleware(), handlers.DeleteWalletUser)

	//Transactions Management
	router.POST("/add-income", middleware.AuthMiddleware(), func(context *gin.Context) { handlers.CreateTransaction(context, true) })
	router.POST("/add-expense", middleware.AuthMiddleware(), func(context *gin.Context) { handlers.CreateTransaction(context, false) })

	router.Run(":8080")
}
