package wallet

import (
	"sync"

	"yyq1025/balance-backend/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

var symbol_cache, decimals_cache sync.Map

func CreateWallet(db *gorm.DB, wallet *Wallet) (int64, error) {
	result := db.Create(wallet)
	return result.RowsAffected, result.Error
}

func QueryWallets(db *gorm.DB, condition *Wallet, wallets *[]Wallet) (int64, error) {
	result := db.Where(condition).Find(wallets)
	return result.RowsAffected, result.Error
}

func DeleteWallets(db *gorm.DB, condition *Wallet) (int64, error) {
	result := db.Where(condition).Delete(&Wallet{})
	return result.RowsAffected, result.Error
}

func GetSymbol(network string, address common.Address, contract *token.Token) (string, error) {
	key := network + ":" + address.String()
	if cache, exist := symbol_cache.Load(key); exist {
		return cache.(string), nil
	}
	symbol, err := contract.Symbol(&bind.CallOpts{})
	if err == nil {
		symbol_cache.Store(key, symbol)
	}
	return symbol, err
}

func GetDecimals(network string, address common.Address, contract *token.Token) (uint8, error) {
	key := network + ":" + address.String()
	if cache, exist := decimals_cache.Load(key); exist {
		return cache.(uint8), nil
	}
	decimals, err := contract.Decimals(&bind.CallOpts{})
	if err == nil {
		decimals_cache.Store(key, decimals)
	}
	return decimals, err
}
