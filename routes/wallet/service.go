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

func addWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, wallet *Wallet) utils.Response {
	balance, err := wallet.getBalance(ctx, rdbCache)
	if err != nil {
		return utils.AddWalletError
	}
	if err := createWallet(ctx, rdbCache, db, wallet); err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}

func deleteBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	if err := deleteWallet(ctx, rdbCache, db, condition); err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"id": condition.ID}}
}

func getBalancesWithPagination(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, p *Pagination) utils.Response {
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
	ch := make(chan Balance)

	for _, wallet := range wallets {
		wg.Add(1)

		go func(w Wallet) {
			defer wg.Done()
			balance, err := w.getBalance(ctx, rdbCache)
			if err != nil {
				log.Print(err)
			}
			ch <- balance
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

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balances": results, "next": p}}
}

func getBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := queryWallet(ctx, rdbCache, db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance, err := wallet.getBalance(ctx, rdbCache)
	if err != nil {
		log.Print(err)
		return utils.GetBalanceError
	}
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}
