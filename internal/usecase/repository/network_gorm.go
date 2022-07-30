package repository

import (
	"context"
	"yyq1025/balance-backend/internal/entity"

	"gorm.io/gorm"
)

type gormNetworkRepository struct {
	db *gorm.DB
}

func NewGormNetworkRepository(db *gorm.DB) entity.NetworkRepository {
	return &gormNetworkRepository{db}
}

func (g *gormNetworkRepository) GetAll(ctx context.Context, networks *[]entity.Network) error {
	return g.db.WithContext(ctx).Find(networks).Error
}
