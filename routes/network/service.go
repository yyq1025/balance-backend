package network

import (
	"log"
	"net/http"

	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetNetWorks(rdbCache *cache.Cache, db *gorm.DB, condition *Network) utils.Response {
	networks := make([]Network, 0)

	if err := QueryNetworks(rdbCache, db, condition, &networks); err != nil {
		log.Print(err)
		return utils.GetNetworkError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"networks": networks}}
}
