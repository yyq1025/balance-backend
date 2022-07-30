package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
	"yyq1025/balance-backend/internal/controller"
	"yyq1025/balance-backend/internal/controller/middleware"
	"yyq1025/balance-backend/internal/usecase"
	"yyq1025/balance-backend/internal/usecase/repository"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	m := middleware.InitMiddleware(jwtValidator, rdb, 3500*time.Millisecond)

	router := gin.Default()
	router.Use(m.CORS())
	router.Use(m.Timeout())
	nr := repository.NewGormNetworkRepository(db)
	ns := usecase.NewNetworkUseCase(nr, rdb)
	controller.NewNetworkHandler(router, ns)

	wr := repository.NewGormWalletRepository(db)
	ws := usecase.NewWalletUseCase(wr, rdb)
	router.Use(m.Auth())
	controller.NewWalletHandler(router, ws)

	log.Fatal(router.Run(":8080"))
}
