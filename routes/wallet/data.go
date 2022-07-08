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

func CreateWallet(rc_cache *cache.Cache, db *gorm.DB, wallet *Wallet) error {
	err := db.Create(wallet).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
	}
	return err
}

func QueryWallets(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet) error {
	var cached_wallet Wallet
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("wallet:%d", condition.ID), &cached_wallet); err == nil {
		*wallets = []Wallet{cached_wallet}
		return nil
	}
	err := db.Where(condition).Find(wallets).Error
	for _, wallet := range *wallets {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: wallet,
			TTL:   time.Hour,
			SetNX: true,
		})
	}
	return err
}

func QueryWallet(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallet *Wallet) error {
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("wallet:%d", condition.ID), wallet); err == nil {
		return nil
	}
	err := db.Where(condition).First(wallet).Error
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("wallet:%d", wallet.ID),
			Value: *wallet,
			TTL:   time.Hour,
		})
	}
	return err
}

func DeleteWallets(rc_cache *cache.Cache, db *gorm.DB, condition *Wallet, wallets *[]Wallet) error {
	err := db.Clauses(clause.Returning{}).Where(condition).Delete(wallets).Error
	for _, wallet := range *wallets {
		rc_cache.Delete(context.TODO(), fmt.Sprintf("wallet:%d", wallet.ID))
	}
	return err
}

func GetSymbol(ctx context.Context, rc_cache *cache.Cache, network string, address common.Address, contract *token.Token) (string, error) {
	var symbol string
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("symbol:%s:%s", network, address.String()), &symbol); err == nil {
		return symbol, nil
	}
	symbol, err := contract.Symbol(&bind.CallOpts{Context: ctx})
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("symbol:%s:%s", network, address.String()),
			Value: symbol,
			TTL:   time.Hour,
		})
	}
	return symbol, err
}

func GetDecimals(ctx context.Context, rc_cache *cache.Cache, network string, address common.Address, contract *token.Token) (uint8, error) {
	var decimals uint8
	if err := rc_cache.Get(context.TODO(), fmt.Sprintf("decimals:%s:%s", network, address.String()), &decimals); err == nil {
		return decimals, nil
	}
	decimals, err := contract.Decimals(&bind.CallOpts{Context: ctx})
	if err == nil {
		rc_cache.Set(&cache.Item{
			Ctx:   context.TODO(),
			Key:   fmt.Sprintf("decimals:%s:%s", network, address.String()),
			Value: decimals,
			TTL:   time.Hour,
		})
	}
	return decimals, err
}
