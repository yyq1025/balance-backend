package wallet

import (
	"context"
	"log"
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
	if err = network.QueryNetwork(rdbCache, db, &network.Network{Name: w.Network}, &walletNetwork); err != nil {
		log.Print(err)
		return
	}
	rpcClient, err := ethclient.Dial(walletNetwork.Url)
	if err != nil {
		log.Print(err)
		return
	}
	if utils.IsZeroAddress(w.Token) {
		b.Symbol = walletNetwork.Symbol
		balance, error := rpcClient.BalanceAt(ctx, w.Address, nil)
		if error != nil {
			log.Print(error)
			err = error
			return
		}
		b.Balance = utils.ToDecimal(balance, 18).String()
		return
	}
	contract, err := token.NewToken(w.Token, rpcClient)
	if err != nil {
		log.Print(err)
		return
	}
	symbol, err := GetSymbol(ctx, rdbCache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		log.Print(err)
		return
	}
	b.Symbol = symbol
	balance, err := contract.BalanceOf(&bind.CallOpts{Context: ctx}, w.Address)
	if err != nil {
		log.Print(err)
		return
	}
	decimals, err := GetDecimals(ctx, rdbCache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		log.Print(err)
		return
	}
	b.Balance = utils.ToDecimal(balance, int(decimals)).String()
	return
}
