package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func createWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, wallet *Wallet) error {
	err := db.WithContext(ctx).Create(wallet).Error
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
	}
	return err
}

func queryWalletsWithPagination(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet, p *Pagination) error {
	if p.IDLte > 0 {
		db = db.Where("id <= ?", p.IDLte)
	}
	err := db.WithContext(ctx).Joins("Network").Where(condition).Order("id desc").Offset(p.Page * p.PageSize).Limit(p.PageSize).Find(wallets).Error
	for _, wallet := range *wallets {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: wallet,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return err
}

func queryWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, wallet *Wallet) error {
	if err := rdbCache.Get(ctx, fmt.Sprintf("wallet:%d", condition.ID), wallet); err == nil {
		return nil
	}
	err := db.WithContext(ctx).Joins("Network").Where(condition).First(wallet).Error
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
	}
	return err
}

func deleteWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) error {
	// User could only delete his own wallet
	err := db.WithContext(ctx).Where(condition).Delete(&Wallet{}).Error
	if err == nil {
		_ = rdbCache.Delete(ctx, fmt.Sprintf("wallet:%d", condition.ID))
	}
	return err
}

// func getTokenBalance(ctx context.Context, rdbCache *cache.Cache, n network.Network, address, token common.Address) (*big.Int, error) {
// 	rpcClient, err := ethclient.Dial(n.URL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if utils.IsZeroAddress(token) {
// 		return rpcClient.BalanceAt(ctx, address, nil)
// 	}
// 	contract, err := erc20.NewErc20(token, rpcClient)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, address)
// }

// func getTokenSymbol(ctx context.Context, rdbCache *cache.Cache, n network.Network, token common.Address) (string, error) {
// 	if utils.IsZeroAddress(token) {
// 		return n.Symbol, nil
// 	}
// 	var symbol string
// 	if err := rdbCache.Get(ctx, fmt.Sprintf("symbol:%s:%s", n.Name, token.String()), &symbol); err == nil {
// 		return symbol, nil
// 	}
// 	rpcClient, err := ethclient.Dial(n.URL)
// 	if err != nil {
// 		return "", err
// 	}
// 	contract, err := erc20.NewErc20(token, rpcClient)
// 	if err != nil {
// 		return "", err
// 	}
// 	symbol, err = contract.Symbol(&bind.CallOpts{Context: ctx})
// 	if err == nil {
// 		_ = rdbCache.Set(&cache.Item{
// 			Ctx:   ctx,
// 			Key:   fmt.Sprintf("symbol:%s:%s", n.Name, token.String()),
// 			Value: symbol,
// 			TTL:   time.Hour,
// 		})
// 	}
// 	return symbol, err
// }

// func getTokenDecimals(ctx context.Context, rdbCache *cache.Cache, n network.Network, token common.Address) (uint8, error) {
// 	if utils.IsZeroAddress(token) {
// 		return 18, nil
// 	}
// 	var decimals uint8
// 	if err := rdbCache.Get(ctx, fmt.Sprintf("decimals:%s:%s", n.Name, token.String()), &decimals); err == nil {
// 		return decimals, nil
// 	}
// 	rpcClient, err := ethclient.Dial(n.URL)
// 	if err != nil {
// 		return 0, err
// 	}
// 	contract, err := erc20.NewErc20(token, rpcClient)
// 	if err != nil {
// 		return 0, err
// 	}
// 	decimals, err = contract.Decimals(&bind.CallOpts{Context: ctx})
// 	if err == nil {
// 		_ = rdbCache.Set(&cache.Item{
// 			Ctx:   ctx,
// 			Key:   fmt.Sprintf("decimals:%s:%s", n.Name, token.String()),
// 			Value: decimals,
// 			TTL:   time.Hour,
// 		})
// 	}
// 	return decimals, err
// }
