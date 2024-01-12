package repository_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/caarlos0/env/v8"
// 	"github.com/yyq1025/balance-backend/config"
// 	"github.com/yyq1025/balance-backend/internal/entity"
// 	"github.com/yyq1025/balance-backend/internal/usecase/repository"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/go-redis/cache/v9"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func wallet(t *testing.T) entity.WalletRepository {
// 	t.Helper()

// 	cfg := &config.Config{}
// 	if err := env.Parse(cfg); err != nil {
// 		t.Fatal(err)
// 	}

// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
// 	})
// 	_, err = rdb.Ping(context.Background()).Result()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rdb.FlushDB(context.Background())

// 	rdbCache := cache.New(&cache.Options{
// 		Redis: rdb,
// 	})

// 	return repository.NewWalletRepository(db, rdbCache)
// }

// func TestWallet(t *testing.T) {
// 	repo := wallet(t)

// 	// Get not cached wallet
// 	var wallet entity.Wallet
// 	if err := repo.GetOne(context.Background(), "1", 1, &wallet); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, entity.Wallet{ID: 1, UserID: "1", NetworkName: "Ethereum", Network: entity.Network{
// 		ChainID:  "0x1",
// 		Name:     "Ethereum",
// 		URL:      "https://eth.public-rpc.com",
// 		Symbol:   "ETH",
// 		Explorer: "https://etherscan.io",
// 	}}, wallet)

// 	// Add wallet
// 	wallet1 := entity.Wallet{
// 		UserID:      "1",
// 		NetworkName: "BSC",
// 		Network: entity.Network{
// 			ChainID:  "0x38",
// 			Name:     "BSC",
// 			URL:      "https://bsc-dataseed.binance.org/",
// 			Symbol:   "BNB",
// 			Explorer: "https://bscscan.com",
// 		},
// 	}
// 	if err := repo.AddOne(context.Background(), &wallet1); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, 2, wallet1.ID)

// 	wallet2 := entity.Wallet{
// 		UserID:      "1",
// 		NetworkName: "Ethereum",
// 		Token:       common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
// 		Network: entity.Network{
// 			ChainID:  "0x1",
// 			Name:     "Ethereum",
// 			URL:      "https://eth.public-rpc.com",
// 			Symbol:   "ETH",
// 			Explorer: "https://etherscan.io",
// 		},
// 	}
// 	if err := repo.AddOne(context.Background(), &wallet2); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, 3, wallet2.ID)

// 	// Add repeat wallet
// 	if err := repo.AddOne(context.Background(), &wallet2); err == nil {
// 		t.Error("should error")
// 	}

// 	// Add wallet timeout
// 	wallet3 := entity.Wallet{
// 		UserID:      "1",
// 		NetworkName: "BSC",
// 	}
// 	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
// 	defer cancel()
// 	if err := repo.AddOne(ctx, &wallet3); err == nil {
// 		t.Error("should error")
// 	}

// 	// Get cached wallet
// 	var wallet4 entity.Wallet
// 	if err := repo.GetOne(context.Background(), "1", 3, &wallet4); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, wallet2, wallet4)

// 	// Get one not found
// 	var wallet5 entity.Wallet
// 	if err := repo.GetOne(context.Background(), "1", 4, &wallet5); err == nil {
// 		t.Error("should error")
// 	}

// 	// Get one timeout
// 	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
// 	defer cancel()
// 	if err := repo.GetOne(ctx, "1", 3, &wallet5); err == nil {
// 		t.Error("should error")
// 	}

// 	// Get many with IDLte
// 	var wallets []entity.Wallet
// 	if err := repo.GetManyWithPagination(context.Background(), "1", &entity.Pagination{IDLte: 2, PageSize: 1}, &wallets); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, []entity.Wallet{wallet1}, wallets)

// 	// Get many without IDLte
// 	var wallets1 []entity.Wallet
// 	if err := repo.GetManyWithPagination(context.Background(), "1", &entity.Pagination{PageSize: 1}, &wallets1); err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, []entity.Wallet{wallet2}, wallets1)

// 	// Get many timeout
// 	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
// 	defer cancel()
// 	if err := repo.GetManyWithPagination(ctx, "1", &entity.Pagination{PageSize: 1}, &wallets1); err == nil {
// 		t.Error("should error")
// 	}

// 	// Delete wallet
// 	if err := repo.DeleteOne(context.Background(), "1", 3); err != nil {
// 		t.Error(err)
// 	}

// 	// Delete wallet timeout
// 	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
// 	defer cancel()
// 	if err := repo.DeleteOne(ctx, "1", 2); err == nil {
// 		t.Error("should error")
// 	}
// }
