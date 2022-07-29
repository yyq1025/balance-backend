package network

import (
	"context"
	"testing"
	"time"
	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestQueryNetworks(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
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

func TestQueryNetworksNoCache(t *testing.T) {
	db := utils.GetDB()
	rdbCache := cache.New(&cache.Options{})
	actual := make([]Network, 0)

	if err := queryAllNetworks(context.Background(), rdbCache, db, &actual); err != nil {
		t.Error(err)
	}
}

func TestQueryNetworksNoDB(t *testing.T) {
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := make([]Network, 0)

	if err := queryAllNetworks(context.Background(), rdbCache, &gorm.DB{}, &actual); err == nil {
		t.Error("expected error")
	}
}

func TestQueryNetworksTimeout(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := make([]Network, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	if err := queryAllNetworks(ctx, rdbCache, db, &actual); err == nil {
		t.Error("expected error")
	}

	if err := queryAllNetworks(ctx, rdbCache, db, &actual); err == nil {
		t.Error("expected error")
	}

	assert.Equal(t, rdbCache.Stats().Hits, uint64(0))
}
