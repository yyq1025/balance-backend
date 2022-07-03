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

func DeleteBalances(db *gorm.DB, condition *Wallet) utils.Response {
	var wallets []Wallet
	rowsAffected, err := DeleteWallets(db, condition, &wallets)
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	if rowsAffected == 0 {
		return utils.FindWalletError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"wallets": wallets}}
}

func GetBalances(db *gorm.DB, condition *Wallet) utils.Response {
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
		return results[i].ID < results[j].ID
	})
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balances": results}}
}

func GetBalance(db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := QueryWallet(db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance := getBalance(db, wallet)
	balance.Wallet = wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}
