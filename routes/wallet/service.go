package wallet

import (
	"context"
	"log"
	"net/http"
	"sync"

	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func AddWallet(ctx context.Context, rdb_cache *cache.Cache, db *gorm.DB, wallet *Wallet) utils.Response {
	balance, err := getBalance(ctx, rdb_cache, db, *wallet)
	if err != nil {
		return utils.AddWalletError
	}
	if err := CreateWallet(rdb_cache, db, wallet); err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	balance.Wallet = *wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}

func DeleteBalances(rdb_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
	err := DeleteWallets(rdb_cache, db, condition, &wallets)
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

func GetBalances(ctx context.Context, rdb_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	wallets := make([]Wallet, 0)
	err := QueryWallets(rdb_cache, db, condition, &wallets)
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
			msg, err := getBalance(ctx, rdb_cache, db, w)
			if err != nil {
				log.Print(err)
			}
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

func GetBalance(ctx context.Context, rdb_cache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := QueryWallet(rdb_cache, db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance, err := getBalance(ctx, rdb_cache, db, wallet)
	if err != nil {
		log.Print(err)
		return utils.GetBalanceError
	}
	balance.Wallet = wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}
