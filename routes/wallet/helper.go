package wallet

import (
	"context"
	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func getBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, w Wallet) (b Balance, err error) {
	var n network.Network
	if err = network.QueryNetwork(ctx, rdbCache, db, &network.Network{Name: w.Network}, &n); err != nil {
		return
	}
	balance, err := getTokenBalance(ctx, rdbCache, n, w.Address, w.Token)
	if err != nil {
		return
	}
	symbol, err := getTokenSymbol(ctx, rdbCache, n, w.Token)
	if err != nil {
		return
	}
	decimals, err := getTokenDecimals(ctx, rdbCache, n, w.Token)
	if err != nil {
		return
	}
	b.Balance = utils.ToDecimal(balance, int(decimals)).InexactFloat64()
	b.Symbol = symbol
	return
}
