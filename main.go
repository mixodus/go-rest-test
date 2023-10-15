package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mixodus/go-rest-test/controllers/bankcontroller"
	"github.com/mixodus/go-rest-test/controllers/playercontroller"
	"github.com/mixodus/go-rest-test/controllers/transactioncontroller"
	"github.com/mixodus/go-rest-test/middleware"
	"github.com/mixodus/go-rest-test/models"
	"github.com/mixodus/go-rest-test/services"
)

func main() {
	godotenv.Load()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// Set the Content-Type header to "application/json"
		c.Header("Content-Type", "application/json")

		// Continue processing other middleware and routes
		c.Next()
	})

	trustedProxies := []string{"localhost", "127.0.0.1"}
	r.SetTrustedProxies(trustedProxies)
	r.Static("/uploads", "./uploads")

	//START DATABASE
	models.ConnectDatabse()

	//START REDIS
	services.InitializeRedis()
	redis := services.GetRedisClient()

	//TEST REDIS
	c := &gin.Context{}
	redis.Set(c, "redis", "REDIS IS WORKING!", time.Hour*1).Err()
	fmt.Print("REDIS VALUE: ")
	fmt.Print(redis.Get(c, "redis").Val())
	fmt.Print("\n\n")

	// ==== API ROUTES ====
	authorized := r.Group("/api", middleware.Authenticate) //middleware
	{
		// PLAYER ROUTES
		authorized.GET("/profile", playercontroller.Profile)
		authorized.DELETE("/logout", playercontroller.Logout)

		//BANK ROUTEs
		authorized.GET("/player/bank", bankcontroller.GetPlayerBank)
		authorized.POST("/player/bank", bankcontroller.AddPlayerBank)
		authorized.DELETE("/player/bank", bankcontroller.RemovePlayerBank)

		//TRANSACTION ROUTES
		authorized.POST("/transaction/topup", transactioncontroller.TopUp) // create transaction topup (debit)
		authorized.POST("/transaction/spent", transactioncontroller.Spent) // create transaction spent (credit)

		//WALLET ROUTES
		authorized.GET("/wallet", transactioncontroller.GetAndUpdateWallet)
	}

	unauth := r.Group("/api")
	{
		// LOGIN & REGISTER ROUTE
		unauth.POST("/register", playercontroller.Register)
		unauth.POST("/login", playercontroller.Login)

		//ADMIN ROUTES (ADMIN ONLY) <- but for now, everyone can access this route because we don't have admin role yet
		unauth.GET("/players", playercontroller.Index)
		unauth.GET("/player/:id", playercontroller.GetPlayerById)
		unauth.PUT("/transaction/debit-success", transactioncontroller.SetAllDebitStatusSuccess)   //set all debit transaction to success
		unauth.PUT("/transaction/credit-success", transactioncontroller.SetAllCreditStatusSuccess) //set all credit transaction to success

		//BANK ROUTES
		unauth.GET("/banks", bankcontroller.BankList)

		//GET IMAGE
		unauth.GET("/image", services.GetImage)
	}

	r.Run(":8080")
}
