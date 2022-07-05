package wallet

import (
	"log"
	"net/http"
	"sync"

	"yyq1025/balance-backend/utils"

	"gorm.io/gorm"
)

func AddWallet(db *gorm.DB, wallet *Wallet) utils.Response {
	balance := getBalance(db, *wallet)
	if balance.Balance == "" {
		return utils.AddWalletError
	}
	rowsAffected, err := CreateWallet(db, wallet)
	if err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	if rowsAffected == 0 {
		return utils.AddWalletError
	}
	balance.Wallet = *wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}

func DeleteBalances(db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
	rowsAffected, err := DeleteWallets(db, condition, &wallets)
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	if rowsAffected == 0 {
		return utils.FindWalletError
	}
	ids := make([]int, 0)
	for _, wallet := range wallets {
		ids = append(ids, wallet.ID)
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"ids": ids}}
}

func GetBalances(db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
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

	results := make([]Balance, 0)
	for result := range ch {
		results = append(results, result)
	}

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
