package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/yyq1025/balance-backend/config"
	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/internal/usecase/repository"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func network(t *testing.T) entity.NetworkRepository {
	t.Helper()

	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(sqlserver.Open(cfg.DB.Dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.EndPoint,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		t.Fatal(err)
	}

	rdb.FlushDB(context.Background())

	rdbCache := cache.New(&cache.Options{
		Redis: rdb,
	})

	return repository.NewNetworkRepository(db, rdbCache)
}

func TestGetAll(t *testing.T) {
	repo := network(t)

	expect := []entity.Network{
		{
			ChainID:  "0x38",
			Name:     "BSC",
			URL:      "https://bsc-dataseed.binance.org/",
			Symbol:   "BNB",
			Explorer: "https://bscscan.com",
		},
		{
			ChainID:  "0x1",
			Name:     "Ethereum",
			URL:      "https://eth.public-rpc.com",
			Symbol:   "ETH",
			Explorer: "https://etherscan.io",
		},
	}

	var networks []entity.Network
	if err := repo.GetAll(context.Background(), &networks); err != nil {
		t.Error(err)
	}

	require.Equal(t, expect, networks)

	// cache
	var networks1 []entity.Network
	if err := repo.GetAll(context.Background(), &networks1); err != nil {
		t.Error(err)
	}

	require.Equal(t, expect, networks1)

	// error
	var networks2 []entity.Network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()

	if err := repo.GetAll(ctx, &networks2); err == nil {
		t.Error("expected error")
	}
}
