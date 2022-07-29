package network

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func queryAllNetworks(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, networks *[]Network) error {
	if err := rdbCache.Get(ctx, "networks", networks); err == nil {
		return nil
	}
	err := db.WithContext(ctx).Order("name asc").Find(networks).Error
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   "networks",
			Value: *networks,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return err
}
