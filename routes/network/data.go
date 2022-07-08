package network

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

// network_cache name:Network
// var network_cache sync.Map
var ctx = context.TODO()

func QueryNetworks(rc_cache *cache.Cache, db *gorm.DB, condition *Network, networks *[]Network) error {
	var cached_network Network
	if err := rc_cache.Get(ctx, fmt.Sprintf("network:%s", condition.Name), &cached_network); err == nil {
		*networks = []Network{cached_network}
		return nil
	}
	err := db.Where(condition).Find(networks).Error
	for _, network := range *networks {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: network,
			TTL:   time.Hour,
			SetNX: true,
		})
		// network_cache.LoadOrStore(network.Name, network)
	}
	return err
}

func QueryNetwork(rc_cache *cache.Cache, db *gorm.DB, condition *Network, network *Network) error {
	if err := rc_cache.Get(ctx, fmt.Sprintf("network:%s", condition.Name), network); err == nil {
		// *network = cached_network.(Network)
		return nil
	}
	err := db.Where(condition).First(network).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("network:%s", network.Name),
			Value: *network,
			TTL:   time.Hour,
		})
		// network_cache.Store(network.Name, *network)
	}
	return err
}
