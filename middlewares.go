package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis_rate/v9"
	"gorm.io/gorm"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Authorization, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, PATCH, DELETE, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func authMiddleware(jwtValidator *validator.Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "empty header"})
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "incorrect header format"})
			return
		}
		claims, err := jwtValidator.ValidateToken(context.Background(), parts[1])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
			return
		}
		c.Set("userID", claims.(*validator.ValidatedClaims).RegisteredClaims.Subject)
		c.Next()
	}
}

func rateLimitMiddleware(limiter *redis_rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		res, err := limiter.Allow(c, userID, redis_rate.PerMinute(10))
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

func dataMiddleware(rdbCache *cache.Cache, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("rdbCache", rdbCache)
		c.Set("db", db)
		c.Next()
	}
}
