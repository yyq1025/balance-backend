package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/yyq1025/balance-backend/config"
	"github.com/yyq1025/balance-backend/internal/controller"
	"github.com/yyq1025/balance-backend/internal/controller/middleware"
	"github.com/yyq1025/balance-backend/internal/usecase"
	"github.com/yyq1025/balance-backend/internal/usecase/ethapi"
	"github.com/yyq1025/balance-backend/internal/usecase/repository"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Username: cfg.Redis.User,
		Password: cfg.Redis.Password,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	rdbCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

	issuerURL, err := url.Parse("https://" + cfg.Auth0.Domain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{cfg.Auth0.Aud},
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	m := middleware.InitMiddleware(jwtValidator, rdb, cfg.Timeout)

	router := gin.Default()
	router.Use(m.CORS())
	router.Use(m.Timeout())
	nr := repository.NewNetworkRepository(db, rdbCache)
	ns := usecase.NewNetworkUseCase(nr)
	controller.NewNetworkHandler(router, ns)

	wr := repository.NewWalletRepository(db, rdbCache)
	we := ethapi.NewWalletEthAPI(rdbCache)
	ws := usecase.NewWalletUseCase(wr, we)
	router.Use(m.Auth())
	controller.NewWalletHandler(router, ws)

	log.Fatal(router.Run(":" + cfg.Port))
}
