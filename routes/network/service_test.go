package network

import (
	"context"
	"net/http"
	"testing"
	"time"
	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"github.com/stretchr/testify/assert"
)

func TestGetAllNetworks(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := getAllNetWorks(context.Background(), rdbCache, db)

	assert.Equal(t, http.StatusOK, actual.Code)
	assert.Greater(t, len(actual.Data["networks"].([]Network)), 0)
}

func TestGetAllNetworksNoCache(t *testing.T) {
	db := utils.GetDB()

	actual := getAllNetWorks(context.Background(), nil, db)

	assert.Equal(t, http.StatusOK, actual.Code)
	assert.Greater(t, len(actual.Data["networks"].([]Network)), 0)
}

func TestGetAllNetworksNoDB(t *testing.T) {
	rdb := utils.GetRedis()
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	actual := getAllNetWorks(context.Background(), rdbCache, nil)

	assert.Equal(t, utils.GetNetworkError, actual)
}

func TestGetAllNetworksTimeout(t *testing.T) {
	db := utils.GetDB()
	rdb := utils.GetRedis()
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	actual := getAllNetWorks(ctx, rdbCache, db)

	assert.Equal(t, utils.GetNetworkError, actual)
}
