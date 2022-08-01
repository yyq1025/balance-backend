package repository

import (
	"context"
	"time"

	"github.com/yyq1025/balance-backend/internal/entity"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

type NetworkRepository struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewNetworkRepository(db *gorm.DB, c *cache.Cache) entity.NetworkRepository {
	return &NetworkRepository{db: db, cache: c}
}

func (n *NetworkRepository) GetAll(ctx context.Context, networks *[]entity.Network) (err error) {
	if err = n.cache.Get(ctx, "networks", networks); err != nil {
		if err = n.db.WithContext(ctx).Order("name asc").Find(networks).Error; err != nil {
			return
		}
		_ = n.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   "networks",
			Value: *networks,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return
}
