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

func AddWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, wallet *Wallet) utils.Response {
	balance, err := getBalance(ctx, rdbCache, db, *wallet)
	if err != nil {
		return utils.AddWalletError
	}
	if err := createWallet(ctx, rdbCache, db, wallet); err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": Result{*wallet, balance}}}
}

func DeleteBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	if err := deleteWallet(ctx, rdbCache, db, condition); err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"id": condition.ID}}
}

func GetBalancesWithPagination(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, p *Pagination) utils.Response {
	wallets := make([]Wallet, 0)
	if err := queryWalletsWithPagination(ctx, rdbCache, db, condition, &wallets, p); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}

	if len(wallets) > 0 && p.IDLte == 0 {
		p.IDLte = wallets[0].ID
	}
	if len(wallets) == p.PageSize {
		p.Page++
	} else {
		p.Page = -1
	}

	var wg sync.WaitGroup
	ch := make(chan Result)

	for _, wallet := range wallets {
		wg.Add(1)

		go func(w Wallet) {
			defer wg.Done()
			balance, err := getBalance(ctx, rdbCache, db, w)
			if err != nil {
				log.Print(err)
			}
			ch <- Result{w, balance}
		}(wallet)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]Result, 0)
	for result := range ch {
		results = append(results, result)
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balances": results, "next": p}}
}

func GetBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := queryWallet(ctx, rdbCache, db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance, err := getBalance(ctx, rdbCache, db, wallet)
	if err != nil {
		log.Print(err)
		return utils.GetBalanceError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": Result{wallet, balance}}}
}
