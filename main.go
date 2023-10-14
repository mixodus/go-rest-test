package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mixodus/go-rest-test/controllers/bankcontroller"
	"github.com/mixodus/go-rest-test/controllers/playercontroller"
	"github.com/mixodus/go-rest-test/middleware"
	"github.com/mixodus/go-rest-test/models"
	"github.com/mixodus/go-rest-test/services"
)

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// Set the Content-Type header to "application/json"
		c.Header("Content-Type", "application/json")

		// Continue processing other middleware and routes
		c.Next()
	})

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
	authorized := r.Group("/api", middleware.Authenticate)
	{
		// PLAYER ROUTES
		authorized.GET("/players", playercontroller.Index)
		authorized.GET("/profile", playercontroller.Profile)
		authorized.DELETE("/logout", playercontroller.Logout)

		//BANK ROUTEs
		authorized.GET("/player/bank", bankcontroller.GetPlayerBank)
		authorized.POST("/player/bank", bankcontroller.AddPlayerBank)
		authorized.DELETE("/player/bank", bankcontroller.RemovePlayerBank)
	}

	unauth := r.Group("/api")
	{
		// LOGIN & REGISTER ROUTE
		unauth.POST("/register", playercontroller.Register)
		unauth.POST("/login", playercontroller.Login)

		//BANK ROUTES
		unauth.GET("/banks", bankcontroller.BankList)
	}

	r.Run()
}
