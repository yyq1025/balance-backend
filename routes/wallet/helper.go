package wallet

import (
	"context"
	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/token"
	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func getBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, w Wallet) (b Balance, err error) {
	var walletNetwork network.Network
	if err = network.QueryNetwork(ctx, rdbCache, db, &network.Network{Name: w.Network}, &walletNetwork); err != nil {
		return
	}
	rpcClient, err := ethclient.Dial(walletNetwork.URL)
	if err != nil {
		return
	}
	if utils.IsZeroAddress(w.Token) {
		b.Symbol = walletNetwork.Symbol
		balance, error := rpcClient.BalanceAt(ctx, w.Address, nil)
		if error != nil {
			b.Balance = -1
			err = error
			return
		}
		b.Balance = utils.ToDecimal(balance, 18).InexactFloat64()
		return
	}
	contract, err := token.NewToken(w.Token, rpcClient)
	if err != nil {
		return
	}
	symbol, err := GetSymbol(ctx, rdbCache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		return
	}
	b.Symbol = symbol
	balance, err := contract.BalanceOf(&bind.CallOpts{Context: ctx}, w.Address)
	if err != nil {
		b.Balance = -1
		return
	}
	decimals, err := GetDecimals(ctx, rdbCache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		b.Balance = -1
		return
	}
	b.Balance = utils.ToDecimal(balance, int(decimals)).InexactFloat64()
	return
}
