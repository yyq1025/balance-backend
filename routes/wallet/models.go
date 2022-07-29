package wallet

import (
	"context"
	"fmt"
	"math/big"
	"time"
	"yyq1025/balance-backend/erc20"
	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/cache/v8"
)

type Wallet struct {
	ID          int             `json:"id"`
	UserID      string          `json:"-"`
	Address     common.Address  `json:"address"`
	NetworkName string          `json:"networkName"`
	Token       common.Address  `json:"token"`
	Network     network.Network `gorm:"foreignKey:NetworkName" json:"network"`
}

func (w Wallet) getTokenBalance(ctx context.Context) (*big.Int, error) {
	rpcClient, err := ethclient.Dial(w.Network.URL)
	if err != nil {
		return nil, err
	}
	if utils.IsZeroAddress(w.Token) {
		return rpcClient.BalanceAt(ctx, w.Address, nil)
	}
	contract, err := erc20.NewErc20(w.Token, rpcClient)
	if err != nil {
		return nil, err
	}
	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, w.Address)
}

func (w Wallet) getTokenSymbol(ctx context.Context, rdbCache *cache.Cache) (string, error) {
	if utils.IsZeroAddress(w.Token) {
		return w.Network.Symbol, nil
	}
	var symbol string
	if err := rdbCache.Get(ctx, fmt.Sprintf("symbol:%s:%s", w.Network.Name, w.Token.String()), &symbol); err == nil {
		return symbol, nil
	}
	rpcClient, err := ethclient.Dial(w.Network.URL)
	if err != nil {
		return "", err
	}
	contract, err := erc20.NewErc20(w.Token, rpcClient)
	if err != nil {
		return "", err
	}
	symbol, err = contract.Symbol(&bind.CallOpts{Context: ctx})
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("symbol:%s:%s", w.Network.Name, w.Token.String()),
			Value: symbol,
			TTL:   time.Hour,
		})
	}
	return symbol, err
}

func (w Wallet) getTokenDecimals(ctx context.Context, rdbCache *cache.Cache) (uint8, error) {
	if utils.IsZeroAddress(w.Token) {
		return 18, nil
	}
	var decimals uint8
	if err := rdbCache.Get(ctx, fmt.Sprintf("decimals:%s:%s", w.Network.Name, w.Token.String()), &decimals); err == nil {
		return decimals, nil
	}
	rpcClient, err := ethclient.Dial(w.Network.URL)
	if err != nil {
		return 0, err
	}
	contract, err := erc20.NewErc20(w.Token, rpcClient)
	if err != nil {
		return 0, err
	}
	decimals, err = contract.Decimals(&bind.CallOpts{Context: ctx})
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("decimals:%s:%s", w.Network.Name, w.Token.String()),
			Value: decimals,
			TTL:   time.Hour,
		})
	}
	return decimals, err
}

func (w Wallet) getBalance(ctx context.Context, rdbCache *cache.Cache) (b Balance, err error) {
	balance, err := w.getTokenBalance(ctx)
	if err != nil {
		return
	}
	symbol, err := w.getTokenSymbol(ctx, rdbCache)
	if err != nil {
		return
	}
	decimals, err := w.getTokenDecimals(ctx, rdbCache)
	if err != nil {
		return
	}
	b.Balance = utils.ToDecimal(balance, int(decimals)).InexactFloat64()
	b.Symbol = symbol
	return
}

type Balance struct {
	Symbol  string  `json:"symbol"`
	Balance float64 `json:"balance"`
}

type Result struct {
	Wallet
	Balance
}

type Pagination struct {
	IDLte    int `json:"idLte"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
