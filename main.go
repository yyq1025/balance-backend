package main

import (
	"log"
	"math/rand"
	"time"

	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/routes/user"
	"yyq1025/balance-backend/routes/wallet"
	"yyq1025/balance-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	db := utils.GetDB()
	rc := utils.GetRedis()
	limiter := redis_rate.NewLimiter(rc)
	sender := utils.NewSender()
	router := gin.Default()
	router.Use(cors.AllowAll())

	user_group := router.Group("/user")
	user_group.Use(dbMiddleware(rc, db))
	{
		user_group.POST("/register", user.RegisterHandler)
		user_group.POST("/code", senderMiddleware(sender), user.SendCodeHandler)
		user_group.POST("/login", user.LoginHandler)
		user_group.PUT("/password", user.ChangePasswordHandler)
	}

	network_group := router.Group("/network")
	network_group.Use(dbMiddleware(rc, db))
	{
		network_group.GET("", network.NetworksHandler)
		network_group.GET("/:network", network.NetworkByNameHandler)
	}

	wallet_group := router.Group("/wallet")
	wallet_group.Use(jwtAuthMiddleware())
	wallet_group.Use(jwtRateLimitMiddleware(limiter))
	wallet_group.Use(dbMiddleware(rc, db))
	{
		wallet_group.POST("/", wallet.CreateWalletHandler)
		// wallet_group.GET("", wallet.GetWalletsHandler)
		wallet_group.DELETE("/:id", wallet.DeleteWalletsHandler)
		wallet_group.GET("/balance", wallet.GetBalancesHandler)
	}

	router.GET("/", jwtAuthMiddleware(), jwtRateLimitMiddleware(limiter), dbMiddleware(rc, db), wallet.GetBalancesHandler)

	log.Fatal(router.Run(":8080"))
}
