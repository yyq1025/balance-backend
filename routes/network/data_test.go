package network

import (
	"context"
	"testing"
	"yyq1025/balance-backend/utils"

	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/cache/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

func TestQueryNetworksCached(t *testing.T) {
	// expected := []Network{{Name: "Ethereum"}}
	db := utils.GetDB()
	rdb := utils.GetRedis()
	// rdb, rdbMock := redismock.NewClientMock()
	rdbCache := cache.New(&cache.Options{
		Redis:        rdb,
		StatsEnabled: true})

	// condition := &Network{Name: "Ethereum"}
	actual := make([]Network, 0)

	// val, _ := rdbCache.Marshal(Network{Name: "Ethereum"})
	// rdbMock.ExpectGet(fmt.Sprintf("network:%s", condition.Name)).SetVal(string(val))

	if err := queryAllNetworks(context.Background(), rdbCache, db, &actual); err != nil {
		t.Error(err)
	}

	if err := queryAllNetworks(context.Background(), rdbCache, db, &actual); err != nil {
		t.Error(err)
	}

	// if err := rdbMock.ExpectationsWereMet(); err != nil {
	// 	t.Error(err)
	// }
	assert.Equal(t, rdbCache.Stats().Hits, rdbCache.Stats().Misses)
}
