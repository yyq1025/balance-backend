package network

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func QueryNetworks(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Network, networks *[]Network) error {
	var cachedNetwork Network
	if err := rdbCache.Get(ctx, fmt.Sprintf("network:%s", condition.Name), &cachedNetwork); err == nil {
		*networks = []Network{cachedNetwork}
		return nil
	}
	err := db.WithContext(ctx).Where(condition).Order("name asc").Find(networks).Error
	for _, network := range *networks {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: network,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return err
}

func QueryNetwork(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB, condition *Network, network *Network) error {
	if err := rdbCache.Get(ctx, fmt.Sprintf("network:%s", condition.Name), network); err == nil {
		return nil
	}
	err := db.WithContext(ctx).Where(condition).First(network).Error
	if err == nil {
		_ = rdbCache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: *network,
			TTL:   time.Hour,
		})
	}
	return err
}
