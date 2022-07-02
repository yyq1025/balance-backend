package wallet

import (
	"context"
	"log"
	"yyq1025/balance-backend/routes/network"
	"yyq1025/balance-backend/token"
	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func getBalance(db *gorm.DB, w Wallet) (b Balance) {
	var walletNetworks []network.Network
	rowsAffected, err := network.QueryNetworks(db, &network.Network{Name: w.Network}, &walletNetworks)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	if rowsAffected == 0 {
		b.Balance = ""
		return
	}
	rpcClient, err := ethclient.Dial(walletNetworks[0].Url)
	if err != nil {
		log.Print(err)
		b.Balance = ""
		return
	}
	if utils.IsZeroAddress(w.Token) {
		b.Symbol = walletNetworks[0].Symbol
		balance, err := rpcClient.BalanceAt(context.Background(), w.Address, nil)
		if err != nil {
			log.Print(err)
			b.Balance = ""
			return
		}
		b.Balance = utils.ToDecimal(balance, 18).String()
		return
	} else {
		contract, err := token.NewToken(w.Token, rpcClient)
		if err != nil {
			log.Print(err)
			b.Balance = ""
			return
		}
		symbol, err := GetSymbol(walletNetworks[0].Name, w.Token, contract)
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
		decimals, err := GetDecimals(walletNetworks[0].Name, w.Token, contract)
		if err != nil {
			log.Print(err)
			b.Balance = ""
			return
		}
		b.Balance = utils.ToDecimal(balance, int(decimals)).String()
		return
	}
}