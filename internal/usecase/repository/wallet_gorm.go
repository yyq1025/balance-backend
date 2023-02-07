package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yyq1025/balance-backend/internal/entity"

	"github.com/go-redis/cache/v9"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewWalletRepository(db *gorm.DB, c *cache.Cache) entity.WalletRepository {
	return &WalletRepository{db: db, cache: c}
}

func (w *WalletRepository) AddOne(ctx context.Context, wallet *entity.Wallet) (err error) {
	if err = w.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return
	}
	_ = w.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("wallet:%d", wallet.ID),
		Value: *wallet,
		TTL:   time.Hour,
		SetNX: true,
	})
	return
}

func (w *WalletRepository) GetOne(ctx context.Context, userID string, id int, wallet *entity.Wallet) (err error) {
	if err = w.cache.Get(ctx, fmt.Sprintf("wallet:%d", id), wallet); err != nil {
		if err = w.db.WithContext(ctx).Joins("Network").Where(&entity.Wallet{UserID: userID, ID: id}).First(wallet).Error; err != nil {
			return
		}
		_ = w.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", id),
			Value: *wallet,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return
}

func (w *WalletRepository) GetManyWithPagination(ctx context.Context, userID string, p *entity.Pagination, wallets *[]entity.Wallet) (err error) {
	db := w.db
	if p.IDLte > 0 {
		db = db.Where("id <= ?", p.IDLte)
	}
	if err = db.WithContext(ctx).Joins("Network").Where(&entity.Wallet{UserID: userID}).Order("id desc").Offset(p.Page * p.PageSize).Limit(p.PageSize).Find(wallets).Error; err != nil {
		return
	}
	for _, wallet := range *wallets {
		_ = w.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: wallet,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return
}

func (w *WalletRepository) DeleteOne(ctx context.Context, userID string, id int) (err error) {
	if err = w.db.WithContext(ctx).Where(&entity.Wallet{UserID: userID, ID: id}).Delete(&entity.Wallet{}).Error; err != nil {
		return
	}
	_ = w.cache.Delete(ctx, fmt.Sprintf("wallet:%d", id))
	return
}
