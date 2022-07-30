package repository

import (
	"context"
	"yyq1025/balance-backend/internal/entity"

	"gorm.io/gorm"
)

type gormWalletRepository struct {
	db *gorm.DB
}

func NewGormWalletRepository(db *gorm.DB) entity.WalletRepository {
	return &gormWalletRepository{db}
}

func (g *gormWalletRepository) AddOne(ctx context.Context, wallet *entity.Wallet) error {
	return g.db.WithContext(ctx).Create(wallet).Error
}

func (g *gormWalletRepository) GetOne(ctx context.Context, condition, wallet *entity.Wallet) error {
	return g.db.WithContext(ctx).Joins("Network").Where(condition).First(wallet).Error
}

func (g *gormWalletRepository) GetManyWithPagination(ctx context.Context, condition *entity.Wallet, wallets *[]entity.Wallet, p *entity.Pagination) error {
	if p.IDLte > 0 {
		return g.db.WithContext(ctx).Joins("Network").Where("id <= ?", p.IDLte).Where(condition).Order("id desc").Offset(p.Page * p.PageSize).Limit(p.PageSize).Find(wallets).Error
	}
	return g.db.WithContext(ctx).Joins("Network").Where(condition).Order("id desc").Offset(p.Page * p.PageSize).Limit(p.PageSize).Find(wallets).Error
}

func (g *gormWalletRepository) DeleteOne(ctx context.Context, condition *entity.Wallet) error {
	return g.db.WithContext(ctx).Where(condition).Delete(&entity.Wallet{}).Error
}
