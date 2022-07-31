package usecase

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
	"yyq1025/balance-backend/internal/entity"
	"yyq1025/balance-backend/pkg/erc20"
	"yyq1025/balance-backend/pkg/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type walletUseCase struct {
	walletRepo entity.WalletRepository
	cache      *cache.Cache
}

func NewWalletUseCase(w entity.WalletRepository, rdb *redis.Client) entity.WalletUseCase {
	return &walletUseCase{
		walletRepo: w,
		cache: cache.New(&cache.Options{
			Redis:      rdb,
			LocalCache: cache.NewTinyLFU(10000, time.Minute),
		})}
}

func (w *walletUseCase) getSymbol(ctx context.Context, wallet *entity.Wallet) (string, error) {
	if util.IsZeroAddress(wallet.Token) {
		return wallet.Network.Symbol, nil
	}
	var symbol string
	if err := w.cache.Get(ctx, fmt.Sprintf("symbol:%s:%s", wallet.Network.Name, wallet.Token.String()), &symbol); err == nil {
		return symbol, nil
	}
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return "", err
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return "", err
	}
	symbol, err = contract.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", err
	}
	_ = w.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("symbol:%s:%s", wallet.Network.Name, wallet.Token.String()),
		Value: symbol,
		TTL:   time.Hour,
	})
	return symbol, nil
}

func (w *walletUseCase) getDecimals(ctx context.Context, wallet *entity.Wallet) (uint8, error) {
	if util.IsZeroAddress(wallet.Token) {
		return 18, nil
	}
	var decimals uint8
	if err := w.cache.Get(ctx, fmt.Sprintf("decimals:%s:%s", wallet.Network.Name, wallet.Token.String()), &decimals); err == nil {
		return decimals, nil
	}
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return 0, err
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return 0, err
	}
	decimals, err = contract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	_ = w.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("decimals:%s:%s", wallet.Network.Name, wallet.Token.String()),
		Value: decimals,
		TTL:   time.Hour,
	})
	return decimals, nil
}

func (w *walletUseCase) getBalance(ctx context.Context, wallet *entity.Wallet) (*big.Int, error) {
	rpcClient, err := ethclient.Dial(wallet.Network.URL)
	if err != nil {
		return nil, err
	}
	if util.IsZeroAddress(wallet.Token) {
		return rpcClient.BalanceAt(ctx, wallet.Address, nil)
	}
	contract, err := erc20.NewErc20(wallet.Token, rpcClient)
	if err != nil {
		return nil, err
	}
	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, wallet.Address)
}

func (w *walletUseCase) AddOne(ctx context.Context, wallet *entity.Wallet) (entity.Balance, error) {
	symbol, err := w.getSymbol(ctx, wallet)
	if err != nil {
		return entity.Balance{}, entity.ErrAddWallet
	}
	decimals, err := w.getDecimals(ctx, wallet)
	if err != nil {
		return entity.Balance{}, entity.ErrAddWallet
	}
	balance, err := w.getBalance(ctx, wallet)
	if err != nil {
		return entity.Balance{}, entity.ErrAddWallet
	}
	if err := w.walletRepo.AddOne(ctx, wallet); err != nil {
		return entity.Balance{}, entity.ErrAddWallet
	}
	_ = w.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("wallet:%d", wallet.ID),
		Value: *wallet,
		TTL:   time.Hour})
	return entity.Balance{Wallet: *wallet, Symbol: symbol, Balance: util.ToDecimal(balance, int(decimals)).InexactFloat64()}, nil
}

func (w *walletUseCase) GetOne(ctx context.Context, condition *entity.Wallet) (entity.Balance, error) {
	var wallet entity.Wallet
	if err := w.cache.Get(ctx, fmt.Sprintf("wallet:%d", condition.ID), &wallet); err != nil {
		if err := w.walletRepo.GetOne(ctx, condition, &wallet); err != nil {
			log.Print(err)
			return entity.Balance{}, entity.ErrGetBalance
		}
		_ = w.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   fmt.Sprintf("wallet:%d", condition.ID),
			Value: wallet,
			TTL:   time.Hour})
	}
	symbol, err := w.getSymbol(ctx, &wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrGetBalance
	}
	decimals, err := w.getDecimals(ctx, &wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrGetBalance
	}
	balance, err := w.getBalance(ctx, &wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrGetBalance
	}
	return entity.Balance{Wallet: wallet, Symbol: symbol, Balance: util.ToDecimal(balance, int(decimals)).InexactFloat64()}, nil
}

func (w *walletUseCase) GetManyWithPagination(ctx context.Context, condition *entity.Wallet, pagination *entity.Pagination) ([]entity.Balance, *entity.Pagination, error) {
	var wallets []entity.Wallet
	if err := w.walletRepo.GetManyWithPagination(ctx, condition, &wallets, pagination); err != nil {
		log.Print(err)
		return nil, nil, entity.ErrFindWallet
	}
	if len(wallets) > 0 && pagination.IDLte == 0 {
		pagination.IDLte = wallets[0].ID
	}
	if len(wallets) == pagination.PageSize {
		pagination.Page++
	} else {
		pagination = nil
	}
	var wg sync.WaitGroup
	wg.Add(len(wallets))
	balances := make([]entity.Balance, len(wallets))
	for i, wallet := range wallets {
		go func(i int, wallet entity.Wallet) {
			defer wg.Done()
			symbol, err := w.getSymbol(ctx, &wallet)
			if err != nil {
				log.Print(err)
				balances[i] = entity.Balance{Wallet: wallet, Symbol: "", Balance: -1}
				return
			}
			decimals, err := w.getDecimals(ctx, &wallet)
			if err != nil {
				log.Print(err)
				balances[i] = entity.Balance{Wallet: wallet, Symbol: "", Balance: -1}
				return
			}
			balance, err := w.getBalance(ctx, &wallet)
			if err != nil {
				log.Print(err)
				balances[i] = entity.Balance{Wallet: wallet, Symbol: "", Balance: -1}
				return
			}
			balances[i] = entity.Balance{Wallet: wallet, Symbol: symbol, Balance: util.ToDecimal(balance, int(decimals)).InexactFloat64()}
		}(i, wallet)
	}
	wg.Wait()
	return balances, pagination, nil
}

func (w *walletUseCase) DeleteOne(ctx context.Context, condition *entity.Wallet) error {
	if err := w.walletRepo.DeleteOne(ctx, condition); err != nil {
		log.Print(err)
		return entity.ErrDeleteWallet
	}
	_ = w.cache.Delete(ctx, fmt.Sprintf("wallet:%d", condition.ID))
	return nil
}
