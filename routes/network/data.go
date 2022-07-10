package network

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func QueryNetworks(rdb_cache *cache.Cache, db *gorm.DB, condition *Network, networks *[]Network) error {
	var cached_network Network
	if err := rdb_cache.Get(context.TODO(), fmt.Sprintf("network:%s", condition.Name), &cached_network); err == nil {
		*networks = []Network{cached_network}
		return nil
	}
	err := db.Where(condition).Find(networks).Error
	for _, network := range *networks {
		rdb_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: network,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return err
}

func QueryNetwork(rdb_cache *cache.Cache, db *gorm.DB, condition *Network, network *Network) error {
	if err := rdb_cache.Get(context.TODO(), fmt.Sprintf("network:%s", condition.Name), network); err == nil {
		return nil
	}
	err := db.Where(condition).First(network).Error
	if err == nil {
		rdb_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: *network,
			TTL:   time.Hour,
		})
	}
	return err
}
