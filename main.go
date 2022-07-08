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
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis_rate/v9"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	db := utils.GetDB()
	rc := utils.GetRedis()
	limiter := redis_rate.NewLimiter(rc)
	rc_cache := cache.New(&cache.Options{
		Redis:        rc,
		LocalCache:   cache.NewTinyLFU(1000, time.Minute),
		StatsEnabled: true,
	})
	router := gin.Default()
	router.Use(cors.AllowAll())

	user_group := router.Group("/user")
	user_group.Use(dataMiddleware(rc_cache, db))
	{
		user_group.POST("/register", user.RegisterHandler)
		user_group.POST("/code", user.SendCodeHandler)
		user_group.POST("/login", user.LoginHandler)
		user_group.PUT("/password", user.ChangePasswordHandler)
	}

	network_group := router.Group("/networks")
	network_group.Use(dataMiddleware(rc_cache, db))
	{
		network_group.GET("", network.GetNetworksHandler)
	}

	wallet_group := router.Group("/wallet")
	wallet_group.Use(authMiddleware())
	wallet_group.Use(rateLimitMiddleware(limiter))
	wallet_group.Use(dataMiddleware(rc_cache, db))
	{
		wallet_group.POST("", wallet.CreateWalletHandler)
		wallet_group.DELETE("/:id", wallet.DeleteWalletsHandler)
		wallet_group.GET("/balances", wallet.GetBalancesHandler)
		wallet_group.GET("/balances/:id", wallet.GetBalanceHandler)
	}

	log.Fatal(router.Run(":8080"))
}
