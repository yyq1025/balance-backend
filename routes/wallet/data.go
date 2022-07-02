package wallet

import (
	"sync"

	"yyq1025/balance-backend/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// wallet_cache: id:Wallet
var wallet_cache sync.Map
var symbol_cache, decimals_cache sync.Map

func CreateWallet(db *gorm.DB, wallet *Wallet) (int64, error) {
	result := db.Create(wallet)
	if result.RowsAffected > 0 {
		wallet_cache.Store(wallet.Id, *wallet)
	}
	return result.RowsAffected, result.Error
}

func QueryWallets(db *gorm.DB, condition *Wallet, wallets *[]Wallet) (int64, error) {
	if cached_wallet, exist := wallet_cache.Load(condition.Id); exist {
		*wallets = []Wallet{cached_wallet.(Wallet)}
		return 1, nil
	}
	result := db.Where(condition).Find(wallets)
	for _, wallet := range *wallets {
		wallet_cache.LoadOrStore(wallet.Id, wallet)
	}
	return result.RowsAffected, result.Error
}

func DeleteWallets(db *gorm.DB, condition *Wallet, wallets *[]Wallet) (int64, error) {
	result := db.Clauses(clause.Returning{}).Where(condition).Delete(wallets)
	for _, wallet := range *wallets {
		wallet_cache.Delete(wallet.Id)
	}
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
