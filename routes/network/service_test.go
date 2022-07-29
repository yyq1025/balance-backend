package network

import (
	"context"
	"net/http"
	"testing"
	"time"
	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllNetworks(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := getAllNetWorks(context.Background(), rdbCache, db)

	assert.Equal(t, http.StatusOK, actual.Code)
	assert.Greater(t, len(actual.Data["networks"].([]Network)), 0)
}

func TestGetAllNetworksNoCache(t *testing.T) {
	db := utils.GetDB()
	rdbCache := cache.New(&cache.Options{})

	actual := getAllNetWorks(context.Background(), rdbCache, db)

	assert.Equal(t, http.StatusOK, actual.Code)
	assert.Greater(t, len(actual.Data["networks"].([]Network)), 0)
}

func TestGetAllNetworksNoDB(t *testing.T) {
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := getAllNetWorks(context.Background(), rdbCache, &gorm.DB{})

	assert.Equal(t, utils.GetNetworkError, actual)
}

func TestGetAllNetworksTimeout(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdb.FlushDB(context.Background())
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	actual := getAllNetWorks(ctx, rdbCache, db)

	assert.Equal(t, utils.GetNetworkError, actual)
}
