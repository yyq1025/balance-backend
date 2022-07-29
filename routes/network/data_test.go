package network

import (
	"context"
	"testing"
	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

func TestQueryNetworksCached(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})
	actual := make([]Network, 0)
	if err := queryAllNetworks(context.Background(), rdbCache, db, &actual); err != nil {
		t.Error(err)
	}

	if err := queryAllNetworks(context.Background(), rdbCache, db, &actual); err != nil {
		t.Error(err)
	}
	assert.Equal(t, rdbCache.Stats().Hits, uint64(1))
	assert.Equal(t, rdbCache.Stats().Hits, rdbCache.Stats().Misses)
}
