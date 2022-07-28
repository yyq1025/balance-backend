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
	if err := CreateWallet(ctx, rdbCache, db, wallet); err != nil {
		log.Print(err)
		return utils.AddWalletError
	}
	balance.Wallet = *wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}

func DeleteBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	if err := DeleteWallet(ctx, rdbCache, db, condition); err != nil {
		log.Print(err)
		return utils.DeleteAddressesError
	}
	// if len(wallets) == 0 {
	// 	return utils.FindWalletError
	// }
	// ids := make([]int, 0)
	// for _, wallet := range wallets {
	// 	ids = append(ids, wallet.ID)
	// }
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"id": condition.ID}}
}

func GetBalancesWithPagination(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, p *Pagination) utils.Response {
	wallets := make([]Wallet, 0)
	if err := QueryWalletsWithPagination(ctx, rdbCache, db, condition, &wallets, p); err != nil {
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
			msg, err := getBalance(ctx, rdbCache, db, w)
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

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balances": results, "next": p}}
}

func GetBalance(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) utils.Response {
	var wallet Wallet
	if err := QueryWallet(ctx, rdbCache, db, condition, &wallet); err != nil {
		log.Print(err)
		return utils.FindWalletError
	}
	balance, err := getBalance(ctx, rdbCache, db, wallet)
	if err != nil {
		log.Print(err)
		return utils.GetBalanceError
	}
	balance.Wallet = wallet
	return utils.Response{Code: http.StatusOK, Data: map[string]any{"balance": balance}}
}
