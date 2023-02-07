package ethapi

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/yyq1025/balance-backend/internal/entity"
	"github.com/yyq1025/balance-backend/pkg/erc20"
	"github.com/yyq1025/balance-backend/pkg/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/cache/v9"
)

type WalletEthAPI struct {
	cache *cache.Cache
}

func NewWalletEthAPI(c *cache.Cache) entity.WalletEthAPI {
	return &WalletEthAPI{c}
}

func (e *WalletEthAPI) getBalance(ctx context.Context, wallet entity.Wallet) (*big.Int, error) {
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return nil, err
	}
	if util.IsZeroAddress(wallet.Token) {
		return rpcClient.BalanceAt(ctx, wallet.Address, nil)
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return nil, err
	}
	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, wallet.Address)
}

func (e *WalletEthAPI) getDecimals(ctx context.Context, wallet entity.Wallet) (uint8, error) {
	if util.IsZeroAddress(wallet.Token) {
		return 18, nil
	}
	var decimals uint8
	if err := e.cache.Get(ctx, fmt.Sprintf("decimals:%s:%s", wallet.Network.Name, wallet.Token.String()), &decimals); err == nil {
		return decimals, nil
	}
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return 0, err
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return 0, err
	}
	decimals, err = contract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	_ = e.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("decimals:%s:%s", wallet.Network.Name, wallet.Token.String()),
		Value: decimals,
		TTL:   time.Hour,
	})
	return decimals, nil
}

func (e *WalletEthAPI) GetSymbol(ctx context.Context, wallet entity.Wallet) (string, error) {
	if util.IsZeroAddress(wallet.Token) {
		return wallet.Network.Symbol, nil
	}
	var symbol string
	if err := e.cache.Get(ctx, fmt.Sprintf("symbol:%s:%s", wallet.Network.Name, wallet.Token.String()), &symbol); err == nil {
		return symbol, nil
	}
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return "", err
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return "", err
	}
	symbol, err = contract.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}
	_ = e.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("symbol:%s:%s", wallet.Network.Name, wallet.Token.String()),
		Value: symbol,
		TTL:   time.Hour,
	})
	return symbol, nil
}

func (e *WalletEthAPI) GetBalance(ctx context.Context, wallet entity.Wallet) (float64, error) {
	decimals, err := e.getDecimals(ctx, wallet)
	if err != nil {
		return 0, err
	}
	balance, err := e.getBalance(ctx, wallet)
	if err != nil {
		return 0, err
	}
	return util.ToDecimal(balance, int(decimals)).InexactFloat64(), nil
}
