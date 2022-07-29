package network

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/cache/v8"
	_ "github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

func TestQueryNetworksCached(t *testing.T) {
	expected := []Network{{Name: "Ethereum"}}
	rdb, rdbMock := redismock.NewClientMock()
	rdbCache := cache.New(&cache.Options{
		Redis: rdb})

	condition := &Network{Name: "Ethereum"}
	actual := make([]Network, 0)

	val, _ := rdbCache.Marshal(Network{Name: "Ethereum"})
	rdbMock.ExpectGet(fmt.Sprintf("network:%s", condition.Name)).SetVal(string(val))

	if err := queryNetworks(context.Background(), rdbCache, nil, condition, &actual); err != nil {
		t.Error(err)
	}

	if err := rdbMock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, actual)
}
