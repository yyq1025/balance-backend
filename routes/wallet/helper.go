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

func getBalance(rc_cache *cache.Cache, db *gorm.DB, w Wallet) (b Balance) {
	var walletNetwork network.Network
	if err := network.QueryNetwork(rc_cache, db, &network.Network{Name: w.Network}, &walletNetwork); err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	rpcClient, err := ethclient.Dial(walletNetwork.Url)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	if utils.IsZeroAddress(w.Token) {
		b.Symbol = walletNetwork.Symbol
		balance, err := rpcClient.BalanceAt(context.Background(), w.Address, nil)
		if err != nil {
			log.Print(err)
			b.Balance = ""
			return
		}
		b.Balance = utils.ToDecimal(balance, 18).String()
		return
	}
	contract, err := token.NewToken(w.Token, rpcClient)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	symbol, err := GetSymbol(rc_cache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		log.Print(err)
	} else {
		b.Symbol = symbol
	}
	balance, err := contract.BalanceOf(&bind.CallOpts{}, w.Address)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	decimals, err := GetDecimals(rc_cache, walletNetwork.Name, w.Token, contract)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	b.Balance = utils.ToDecimal(balance, int(decimals)).String()
	return
}
