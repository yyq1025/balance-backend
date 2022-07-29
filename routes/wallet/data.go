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
