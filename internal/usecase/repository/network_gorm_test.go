package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/internal/usecase/repository"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func network(t *testing.T) entity.NetworkRepository {
	t.Helper()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
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
