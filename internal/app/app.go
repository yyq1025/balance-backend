package app

import (
	"context"
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
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Run(cfg *config.Config) {
	db, err := gorm.Open(mysql.Open(cfg.DB.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.EndPoint,
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

	log.Fatal(router.Run())
}
