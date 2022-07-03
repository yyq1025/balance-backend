package network

import (
	"sync"

	"gorm.io/gorm"
)

// network_cache name:Network
var network_cache sync.Map

func QueryNetworks(db *gorm.DB, condition *Network, networks *[]Network) (int64, error) {
	if cached_network, exist := network_cache.Load(condition.Name); exist {
		*networks = []Network{cached_network.(Network)}
		return 1, nil
	}
	result := db.Where(condition).Find(networks)
	for _, network := range *networks {
		network_cache.LoadOrStore(network.Name, network)
	}
	return result.RowsAffected, result.Error
}

func QueryNetwork(db *gorm.DB, condition *Network, network *Network) error {
	if cached_network, exist := network_cache.Load(condition.Name); exist {
		*network = cached_network.(Network)
		return nil
	}
	err := db.Where(condition).First(network).Error
	if err == nil {
		network_cache.Store(network.Name, *network)
	}
	return err
}
