package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	jwtValidator *validator.Validator
	rateLimiter  *redis_rate.Limiter
	timeout      time.Duration
}

func InitMiddleware(jwtValidator *validator.Validator, rdb *redis.Client, timeout time.Duration) *Middleware {
	return &Middleware{
		jwtValidator: jwtValidator,
		rateLimiter:  redis_rate.NewLimiter(rdb),
		timeout:      timeout,
	}
}

func (m *Middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Accept, Accept-Language, Content-Language, X-Requested-With, Cache-Control, Pragma, Date")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (m *Middleware) Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), m.timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *Middleware) Auth() gin.HandlerFunc {
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
		claims, err := m.jwtValidator.ValidateToken(context.Background(), parts[1])
		if err != nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
			return
		}
		userID := claims.(*validator.ValidatedClaims).RegisteredClaims.Subject
		res, err := m.rateLimiter.Allow(c, userID, redis_rate.PerMinute(10))
		if err != nil {
			log.Print(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "rate limit error"})
			return
		}
		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "reach rate limit"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
