package main

import (
	"net/http"
	"strconv"
	"strings"

	"yyq1025/balance-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"gorm.io/gorm"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, PATCH, DELETE, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "empty header"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "incorrect header format"})
			return
		}
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
			return
		}
		c.Set("userId", claims.UserId)
		c.Next()
	}
}

func jwtRateLimitMiddleware(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("userId").(int)
		res, err := limiter.Allow(c, strconv.Itoa(userId), redis_rate.PerMinute(10))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "rate limit error"})
			return
		}
		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "reach rate limit"})
			return
		}
		c.Next()
	}
}

// dbMiddleware will add the db connection to the context
func dbMiddleware(rc *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rc", rc)
		c.Set("db", db)
		c.Next()
	}
}

func senderMiddleware(s *utils.Sender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("sender", s)
		c.Next()
	}
}
