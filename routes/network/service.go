package network

import (
	"log"
	"net/http"

	"yyq1025/balance-backend/utils"

	"gorm.io/gorm"
)

func GetAllNetWorks(db *gorm.DB) utils.Response {
	networks := make([]Network, 0)

	_, err := QueryNetworks(db, &Network{}, &networks)

	if err != nil {
		log.Print(err)
		return utils.GetNetworkError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"networks": networks}}
}
