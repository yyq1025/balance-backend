package wallet

import (
	"log"
	"net/http"
	"sync"
	"time"

	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func AddWallet(rc_cache *cache.Cache, db *gorm.DB, wallet *Wallet) utils.Response {
	balance := getBalance(rc_cache, db, *wallet)
	if balance.Balance == "" {
		return utils.AddWalletError
	}
	if err := CreateWallet(rc_cache, db, wallet); err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	// if rowsAffected == 0 {
	// 	return utils.AddWalletError
	// }
	balance.Wallet = *wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}

func DeleteBalances(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
	err := DeleteWallets(rc_cache, db, condition, &wallets)
	if err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	if len(wallets) == 0 {
		return utils.FindWalletError
	}
	ids := make([]int, 0)
	for _, wallet := range wallets {
		ids = append(ids, wallet.ID)
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"ids": ids}}
}

func GetBalances(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
	err := QueryWallets(rc_cache, db, condition, &wallets)
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
			start := time.Now()
			msg := getBalance(rc_cache, db, w)
			msg.Wallet = w
			log.Print(w, time.Since(start))
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

func GetBalance(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := QueryWallet(rc_cache, db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance := getBalance(rc_cache, db, wallet)
	balance.Wallet = wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}
