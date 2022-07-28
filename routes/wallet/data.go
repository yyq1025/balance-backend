package wallet

import (
	"context"
	"fmt"
	"time"

	"yyq1025/balance-backend/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func CreateWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, wallet *Wallet) error {
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

func QueryWalletsWithPagination(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet, p *Pagination) error {
	if p.IDLte > 0 {
		db = db.Where("id <= ?", p.IDLte)
	}
	err := db.WithContext(ctx).Where(condition).Order("id desc").Offset(p.Page * p.PageSize).Limit(p.PageSize).Find(wallets).Error
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

func QueryWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet, wallet *Wallet) error {
	if err := rdbCache.Get(ctx, fmt.Sprintf("wallet:%d", condition.ID), wallet); err == nil {
		return nil
	}
	err := db.WithContext(ctx).Where(condition).First(wallet).Error
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

func DeleteWallet(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Wallet) error {
	// User could only delete his own wallet
	err := db.WithContext(ctx).Where(condition).Delete(condition).Error
	if err == nil {
		_ = rdbCache.Delete(ctx, fmt.Sprintf("wallet:%d", condition.ID))
	}
	return err
}

func GetSymbol(ctx context.Context, rdbCache *cache.Cache, network string, address common.Address, contract *token.Token) (string, error) {
	var symbol string
	if err := rdbCache.Get(ctx, fmt.Sprintf("symbol:%s:%s", network, address.String()), &symbol); err == nil {
		return symbol, nil
	}
	symbol, err := contract.Symbol(&bind.CallOpts{Context: ctx})
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("symbol:%s:%s", network, address.String()),
			Value: symbol,
			TTL:   time.Hour,
		})
	}
	return symbol, err
}

func GetDecimals(ctx context.Context, rdbCache *cache.Cache, network string, address common.Address, contract *token.Token) (uint8, error) {
	var decimals uint8
	if err := rdbCache.Get(ctx, fmt.Sprintf("decimals:%s:%s", network, address.String()), &decimals); err == nil {
		return decimals, nil
	}
	decimals, err := contract.Decimals(&bind.CallOpts{Context: ctx})
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("decimals:%s:%s", network, address.String()),
			Value: decimals,
			TTL:   time.Hour,
		})
	}
	return decimals, err
}
