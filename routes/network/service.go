package network

import (
	"context"
	"log"
	"net/http"
	"sync"

	"yyq1025/balance-backend/utils"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func GetAllNetWorks(db *gorm.DB) utils.Response {
	var networks []Network

	_, err := QueryNetworks(db, &Network{}, &networks)

	if err != nil {
		log.Print(err)
		return utils.GetNetworkError
	}

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"networks": networks}}
}

func GetNetworkInfoByName(db *gorm.DB, name string) utils.Response {
	var networks []Network

	rowsAffected, err := QueryNetworks(db, &Network{Name: name}, &networks)

	if err != nil {
		log.Print(err)
		return utils.UnsupportNetworkError
	}

	if rowsAffected == 0 {
		return utils.UnsupportNetworkError
	}

	rpcClient, err := ethclient.Dial(networks[0].Url)

	if err != nil {
		log.Print(err)
		return utils.EthError
	}

	info := Info{Network: networks[0]}
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		blockNumber, err := rpcClient.BlockNumber(context.Background())

		if err != nil {
			log.Print(err)
		} else {
			info.BlockNumber = blockNumber
		}

		wg.Done()
	}()

	go func() {
		gasPrice, err := rpcClient.SuggestGasPrice(context.Background())

		if err != nil {
			log.Print(err)
		} else {
			info.GasPrice = utils.ToDecimal(gasPrice, 9).InexactFloat64()
		}

		wg.Done()
	}()

	wg.Wait()

	return utils.Response{Code: http.StatusOK, Data: map[string]any{"info": info}}
}
