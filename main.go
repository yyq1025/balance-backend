package main

import (
	"log"
	"math/rand"
	"time"

	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/routes/wallet"
	"yyq1025/balance-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis_rate/v9"
)

func main() {
	rand.Seed(time.Now().UnixMicro())
	db := utils.GetDB()
	rdb := utils.GetRedis()
	limiter := redis_rate.NewLimiter(rdb)
	rdbCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	jwtValidator := utils.GetValidator()
	router := gin.Default()
	router.Use(corsMiddleware())
	networkGroup := router.Group("/networks")
	networkGroup.Use(dataMiddleware(rdbCache, db))
	{
		networkGroup.GET("", network.GetNetworksHandler)
	}

	walletGroup := router.Group("/wallet")
	walletGroup.Use(authMiddleware(jwtValidator))
	walletGroup.Use(rateLimitMiddleware(limiter))
	walletGroup.Use(dataMiddleware(rdbCache, db))
	{
		walletGroup.POST("", wallet.CreateWalletHandler)
		walletGroup.DELETE("/:id", wallet.DeleteWalletHandler)
		walletGroup.GET("/balances", wallet.GetBalancesHandler)
		walletGroup.GET("/balances/:id", wallet.GetBalanceHandler)
	}

	log.Fatal(router.Run(":8080"))
}
