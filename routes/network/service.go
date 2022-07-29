package network

import (
	"context"
	"log"
	"net/http"

	"yyq1025/balance-backend/utils"

	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func getAllNetWorks(ctx context.Context, rdbCache *cache.Cache, db *gorm.DB) utils.Response {
	networks := make([]Network, 0)

	if err := queryAllNetworks(ctx, rdbCache, db, &networks); err != nil {
		log.Print(err)
		return utils.GetNetworkError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"networks": networks}}
}
