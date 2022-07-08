package wallet

import (
	"context"
	"fmt"
	"time"

	"yyq1025/balance-backend/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// wallet_cache: id:Wallet
// var wallet_cache sync.Map
// var symbol_cache, decimals_cache sync.Map
var ctx = context.TODO()

func CreateWallet(rc_cache *cache.Cache, db *gorm.DB, wallet *Wallet) error {
	err := db.Create(wallet).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
		// wallet_cache.Store(wallet.ID, *wallet)
	}
	return err
}

func QueryWallets(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet) error {
	var cached_wallet Wallet
	if err := rc_cache.Get(ctx, fmt.Sprintf("wallet:%d", condition.ID), &cached_wallet); err == nil {
		*wallets = []Wallet{cached_wallet}
		return nil
	}
	err := db.Where(condition).Find(wallets).Error
	for _, wallet := range *wallets {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: wallet,
			TTL:   time.Hour,
			SetNX: true,
		})
		// wallet_cache.LoadOrStore(wallet.ID, wallet)
	}
	return err
}

func QueryWallet(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallet *Wallet) error {
	if err := rc_cache.Get(ctx, fmt.Sprintf("wallet:%d", condition.ID), wallet); err == nil {
		// *wallet = cached_wallet.(Wallet)
		return nil
	}
	err := db.Where(condition).First(wallet).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
		// wallet_cache.Store(wallet.ID, *wallet)
	}
	return err
}

func DeleteWallets(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet) error {
	err := db.Clauses(clause.Returning{}).Where(condition).Delete(wallets).Error
	for _, wallet := range *wallets {
		rc_cache.Delete(ctx, fmt.Sprintf("wallet:%d", wallet.ID))
		// wallet_cache.Delete(wallet.ID)
	}
	return err
}

func GetSymbol(rc_cache *cache.Cache, network string, address common.Address, contract *token.Token) (string, error) {
	var symbol string
	if err := rc_cache.Get(ctx, fmt.Sprintf("symbol:%s:%s", network, address.String()), &symbol); err == nil {
		return symbol, nil
	}
	// key := network + ":" + address.String()
	// if cache, exist := symbol_cache.Load(key); exist {
	// 	return cache.(string), nil
	// }
	symbol, err := contract.Symbol(&bind.CallOpts{})
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("symbol:%s:%s", network, address.String()),
			Value: symbol,
			TTL:   time.Hour,
		})
		// symbol_cache.Store(key, symbol)
	}
	return symbol, err
}

func GetDecimals(rc_cache *cache.Cache, network string, address common.Address, contract *token.Token) (uint8, error) {
	var decimals uint8
	if err := rc_cache.Get(ctx, fmt.Sprintf("decimals:%s:%s", network, address.String()), &decimals); err == nil {
		return decimals, nil
	}
	// key := network + ":" + address.String()
	// if cache, exist := decimals_cache.Load(key); exist {
	// 	return cache.(uint8), nil
	// }
	decimals, err := contract.Decimals(&bind.CallOpts{})
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("decimals:%s:%s", network, address.String()),
			Value: decimals,
			TTL:   time.Hour,
		})
		// decimals_cache.Store(key, decimals)
	}
	return decimals, err
}
