package ethapi_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/internal/usecase/ethapi"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

var ethChain = entity.Network{
	ChainID:  "0x1",
	Name:     "Ethereum",
	URL:      "https://eth.public-rpc.com",
	Symbol:   "ETH",
	Explorer: "https://etherscan.io",
}

func ethAPI(t *testing.T) entity.WalletEthAPI {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		t.Fatal(err)
	}

	rdb.FlushDB(context.Background())

	rdbCache := cache.New(&cache.Options{
		Redis: rdb,
	})

	return ethapi.NewWalletEthAPI(rdbCache)
}

func TestEthAPI(t *testing.T) {
	api := ethAPI(t)

	// get native token symbol
	symbol, err := api.GetSymbol(context.Background(), entity.Wallet{Network: ethChain})
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "ETH", symbol)

	// get DAI token symbol
	symbol, err = api.GetSymbol(context.Background(), entity.Wallet{Network: ethChain, Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")})
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "DAI", symbol)

	// get cached token symbol
	symbol, err = api.GetSymbol(context.Background(), entity.Wallet{Network: ethChain, Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")})
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "DAI", symbol)

	// get symbol invalid url
	_, err = api.GetSymbol(context.Background(), entity.Wallet{Network: entity.Network{URL: "abcd"}, Token: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")})
	if err == nil {
		t.Error("should be error")
	}

	// get symbol timeout
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	_, err = api.GetSymbol(ctx, entity.Wallet{Network: ethChain, Token: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")})
	if err == nil {
		t.Error("should be error")
	}

	// get native token balance
	balance, err := api.GetBalance(context.Background(), entity.Wallet{Network: ethChain})
	if err != nil {
		t.Error(err)
	}
	require.Greater(t, balance, float64(0))

	// get DAI token balance
	balance, err = api.GetBalance(context.Background(), entity.Wallet{Network: ethChain, Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")})
	if err != nil {
		t.Error(err)
	}
	require.Greater(t, balance, float64(0))

	// get DAI token balance with cached decimals
	balance, err = api.GetBalance(context.Background(), entity.Wallet{Network: ethChain, Token: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F")})
	if err != nil {
		t.Error(err)
	}
	require.Greater(t, balance, float64(0))

	// get native token balance with invalid url
	_, err = api.GetBalance(context.Background(), entity.Wallet{Network: entity.Network{URL: "abcd"}})
	if err == nil {
		t.Error("should be error")
	}

	// get USDT token balance with invalid url
	_, err = api.GetBalance(context.Background(), entity.Wallet{Network: entity.Network{URL: "abcd"}, Token: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")})
	if err == nil {
		t.Error("should be error")
	}

	// get USDT token balance with timeout
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Minute))
	defer cancel()
	_, err = api.GetBalance(ctx, entity.Wallet{Network: ethChain, Token: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")})
	if err == nil {
		t.Error("should be error")
	}
}
