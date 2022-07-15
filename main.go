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
	rdb_cache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	jwtValidator := utils.GetValidator()
	router := gin.Default()
	// router.Use(cors.AllowAll())
	router.Use(corsMiddleware())
	network_group := router.Group("/networks")
	network_group.Use(dataMiddleware(rdb_cache, db))
	{
		network_group.GET("", network.GetNetworksHandler)
	}

	wallet_group := router.Group("/wallet")
	wallet_group.Use(authMiddleware(jwtValidator))
	wallet_group.Use(rateLimitMiddleware(limiter))
	wallet_group.Use(dataMiddleware(rdb_cache, db))
	{
		wallet_group.POST("", wallet.CreateWalletHandler)
		wallet_group.DELETE("/:id", wallet.DeleteWalletsHandler)
		wallet_group.GET("/balances", wallet.GetBalancesHandler)
		wallet_group.GET("/balances/:id", wallet.GetBalanceHandler)
	}

	log.Fatal(router.Run(":8080"))
}
