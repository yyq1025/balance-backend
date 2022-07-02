package wallet

import (
	"log"
	"net/http"
	"sort"
	"sync"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

func AddWallet(db *gorm.DB, userId int, address, network, tokenAddress, tag string) utils.Response {
	wallet := Wallet{
		UserId:  userId,
		Address: common.HexToAddress(address),
		Network: network,
		Token:   common.HexToAddress(tokenAddress),
		Tag:     tag,
	}
	balance := getBalance(db, wallet)
	if balance.Balance == "" {
		return utils.AddWalletError
	}
	rowsAffected, err := CreateWallet(db, &wallet)
	if err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	if rowsAffected == 0 {
		return utils.AddWalletError
	}
	balance.Wallet = wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
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

func DeleteWalletsByIds(db *gorm.DB, condition *Wallet) utils.Response {
	// condition := &Wallet{
	// 	Id:     Id,
	// 	UserId: userId,
	// }

	rowsAffected, err := DeleteWallets(db, condition, &[]Wallet{})
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	if rowsAffected == 0 {
		return utils.FindWalletError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"message": "delete success"}}
}

func GetBalanceByParams(db *gorm.DB, condition *Wallet) utils.Response {
	// condition := &Wallet{
	// 	UserId:  userId,
	// 	Address: common.HexToAddress(address),
	// 	Network: networkName,
	// 	Token:   common.HexToAddress(tokenAddress),
	// 	Tag:     tag,
	// }
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

		go func(w Wallet) {
			defer wg.Done()
			msg := getBalance(db, w)
			msg.Wallet = w
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

	sort.Slice(results, func(i, j int) bool {
		return results[i].Id < results[j].Id
	})
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balances": results}}
}
