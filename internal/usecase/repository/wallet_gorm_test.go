package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/internal/usecase/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func wallet(t *testing.T) entity.WalletRepository {
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

	return repository.NewWalletRepository(db, rdbCache)
}

func TestWallet(t *testing.T) {
	repo := wallet(t)

	// Get not cached wallet
	var wallet entity.Wallet
	if err := repo.GetOne(context.Background(), entity.Wallet{ID: 1, UserID: "1"}, &wallet); err != nil {
		t.Error(err)
	}
	require.Equal(t, entity.Wallet{ID: 1, UserID: "1", NetworkName: "Ethereum", Network: entity.Network{
		ChainID:  "0x1",
		Name:     "Ethereum",
		URL:      "https://eth.public-rpc.com",
		Symbol:   "ETH",
		Explorer: "https://etherscan.io",
	}}, wallet)

	// Add wallet
	wallet1 := entity.Wallet{
		UserID:      "1",
		NetworkName: "BSC",
	}
	if err := repo.AddOne(context.Background(), &wallet1); err != nil {
		t.Error(err)
	}
	require.Equal(t, 2, wallet1.ID)

	wallet2 := entity.Wallet{
		UserID:      "1",
		NetworkName: "Ethereum",
		Token:       common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
	}
	if err := repo.AddOne(context.Background(), &wallet2); err != nil {
		t.Error(err)
	}
	require.Equal(t, 3, wallet2.ID)

	// Add repeat wallet
	if err := repo.AddOne(context.Background(), &wallet2); err == nil {
		t.Error("should error")
	}

	// Add wallet timeout
	wallet3 := entity.Wallet{
		UserID:      "1",
		NetworkName: "BSC",
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	if err := repo.AddOne(ctx, &wallet3); err == nil {
		t.Error("should error")
	}

	// Get cached wallet
	var wallet4 entity.Wallet
	if err := repo.GetOne(context.Background(), entity.Wallet{ID: 3, UserID: "1"}, &wallet4); err != nil {
		t.Error(err)
	}
	require.Equal(t, wallet2, wallet4)

	// Get one not found
	var wallet5 entity.Wallet
	if err := repo.GetOne(context.Background(), entity.Wallet{ID: 4, UserID: "1"}, &wallet5); err == nil {
		t.Error("should error")
	}

	// Get one timeout
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	if err := repo.GetOne(ctx, entity.Wallet{ID: 3, UserID: "1"}, &wallet5); err == nil {
		t.Error("should error")
	}

	// Get many with IDLte
	var wallets []entity.Wallet
	if err := repo.GetManyWithPagination(context.Background(), entity.Wallet{UserID: "1"}, &wallets, &entity.Pagination{IDLte: 2, PageSize: 1}); err != nil {
		t.Error(err)
	}
	require.Equal(t, []entity.Wallet{wallet1}, wallets)

	// Get many without IDLte
	var wallets1 []entity.Wallet
	if err := repo.GetManyWithPagination(context.Background(), entity.Wallet{UserID: "1"}, &wallets1, &entity.Pagination{PageSize: 1}); err != nil {
		t.Error(err)
	}
	require.Equal(t, []entity.Wallet{wallet2}, wallets1)

	// Get many timeout
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	if err := repo.GetManyWithPagination(ctx, entity.Wallet{UserID: "1"}, &wallets1, &entity.Pagination{PageSize: 1}); err == nil {
		t.Error("should error")
	}

	// Delete wallet
	if err := repo.DeleteOne(context.Background(), entity.Wallet{ID: 3, UserID: "1"}); err != nil {
		t.Error(err)
	}

	// Delete wallet timeout
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	if err := repo.DeleteOne(ctx, entity.Wallet{ID: 2, UserID: "1"}); err == nil {
		t.Error("should error")
	}
}
