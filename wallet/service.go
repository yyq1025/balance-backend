package wallet

import (
	"context"
	"log"
	"net/http"
	"sync"

	"yyq1025/balance-backend/network"
	"yyq1025/balance-backend/token"
	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func AddWallet(db *gorm.DB, userId int, address, network, tokenAddress, tag string) utils.Response {
	wallet := &Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: network,
		Token:   common.HexToAddress(tokenAddress),
		Tag:     tag,
	}
	rowsAffected, err := CreateWallet(db, wallet)
	if err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	if rowsAffected == 0 {
		return utils.AddWalletError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"status": "add wallet success"}}
}

func GetWalletsByParams(db *gorm.DB, userId int, address, networkName, tokenAddress, tag string) utils.Response {
	condition := &Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: networkName,
		Token:   common.HexToAddress(tokenAddress),
		Tag:     tag,
	}
	var wallets []Wallet
	_, err := QueryWallets(db, condition, &wallets)
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"wallets": wallets}}
}

func DeleteWalletsByParams(db *gorm.DB, userId int, address, networkName, tokenAddress, tag string) utils.Response {
	condition := &Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: networkName,
		Token:   common.HexToAddress(tokenAddress),
		Tag:     tag,
	}

	rowsAffected, err := DeleteWallets(db, condition)
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	if rowsAffected == 0 {
		return utils.FindWalletError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"status": "delete success"}}
}

func GetBalanceByParams(db *gorm.DB, userId int, address, networkName, tokenAddress, tag string) utils.Response {
	condition := &Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: networkName,
		Token:   common.HexToAddress(tokenAddress),
		Tag:     tag,
	}
	var wallets []Wallet
	_, err := QueryWallets(db, condition, &wallets)
	if err != nil {
		log.Print(err)
		return utils.FindWalletError
	}

	var wg sync.WaitGroup
	ch := make(chan Balance)

	for _, wallet := range wallets {
		wg.Add(1)

		go func(wallet Wallet) {
			defer wg.Done()
			msg := Balance{Wallet: wallet}
			var walletNetworks []network.Network
			rowsAffected, err := network.QueryNetworks(db, &network.Network{Name: wallet.Network}, &walletNetworks)
			if err != nil {
				log.Print(err)
				msg.Balance = "cannot get balance"
				ch <- msg
				return
			}
			if rowsAffected == 0 {
				msg.Balance = "cannot get balance"
				ch <- msg
				return
			}
			rpcClient, err := ethclient.Dial(walletNetworks[0].Url)
			if err != nil {
				log.Print(err)
				msg.Balance = "cannot get balance"
				ch <- msg
				return
			}
			if utils.IsZeroAddress(wallet.Token) {
				msg.Symbol = walletNetworks[0].Symbol
				balance, err := rpcClient.BalanceAt(context.Background(), wallet.Address, nil)
				if err != nil {
					log.Print(err)
					msg.Balance = "cannot get balance"
				} else {
					msg.Balance = utils.ToDecimal(balance, 18)
				}
			} else {
				contract, err := token.NewToken(wallet.Token, rpcClient)
				if err != nil {
					log.Print(err)
					msg.Balance = "cannot get balance"
					ch <- msg
					return
				}
				symbol, err := GetSymbol(walletNetworks[0].Name, wallet.Token, contract)
				if err != nil {
					log.Print(err)
				} else {
					msg.Symbol = symbol
				}
				balance, err := contract.BalanceOf(&bind.CallOpts{}, wallet.Address)
				if err != nil {
					log.Print(err)
					msg.Balance = "cannot get balance"
					ch <- msg
					return
				}
				decimals, err := GetDecimals(walletNetworks[0].Name, wallet.Token, contract)
				if err != nil {
					log.Print(err)
					msg.Balance = "cannot get balance"
					ch <- msg
					return
				}
				msg.Balance = utils.ToDecimal(balance, int(decimals))
			}
			ch <- msg
		}(wallet)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []Balance
	for result := range ch {
		results = append(results, result)
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"wallets": results}}
}
