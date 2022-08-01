package usecase

import (
	"context"
	"log"
	"sync"

	"github.com/yyq1025/balance-backend/internal/entity"
)

type walletUseCase struct {
	repo   entity.WalletRepository
	ethAPI entity.WalletEthAPI
}

func NewWalletUseCase(wr entity.WalletRepository, we entity.WalletEthAPI) entity.WalletUseCase {
	return &walletUseCase{repo: wr, ethAPI: we}
}

func (w *walletUseCase) AddOne(ctx context.Context, wallet *entity.Wallet) (entity.Balance, error) {
	symbol, err := w.ethAPI.GetSymbol(ctx, *wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrAddWallet
	}
	balance, err := w.ethAPI.GetBalance(ctx, *wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrAddWallet
	}
	if err := w.repo.AddOne(ctx, wallet); err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrAddWallet
	}
	return entity.Balance{Wallet: *wallet, Symbol: symbol, Balance: balance}, nil
}

func (w *walletUseCase) GetOne(ctx context.Context, condition entity.Wallet) (entity.Balance, error) {
	var wallet entity.Wallet
	if err := w.repo.GetOne(ctx, condition, &wallet); err != nil {
		return entity.Balance{}, entity.ErrGetBalance
	}
	symbol, err := w.ethAPI.GetSymbol(ctx, wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrGetBalance
	}
	balance, err := w.ethAPI.GetBalance(ctx, wallet)
	if err != nil {
		log.Print(err)
		return entity.Balance{}, entity.ErrGetBalance
	}
	return entity.Balance{Wallet: wallet, Symbol: symbol, Balance: balance}, nil
}

func (w *walletUseCase) GetManyWithPagination(ctx context.Context, condition entity.Wallet, pagination *entity.Pagination) ([]entity.Balance, *entity.Pagination, error) {
	var wallets []entity.Wallet
	if err := w.repo.GetManyWithPagination(ctx, condition, &wallets, pagination); err != nil {
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
			symbol, err := w.ethAPI.GetSymbol(ctx, wallet)
			if err != nil {
				log.Print(err)
				balances[i] = entity.Balance{Wallet: wallet, Symbol: "", Balance: -1}
				return
			}
			balance, err := w.ethAPI.GetBalance(ctx, wallet)
			if err != nil {
				log.Print(err)
				balances[i] = entity.Balance{Wallet: wallet, Symbol: symbol, Balance: -1}
				return
			}
			balances[i] = entity.Balance{Wallet: wallet, Symbol: symbol, Balance: balance}
		}(i, wallet)
	}
	wg.Wait()
	return balances, pagination, nil
}

func (w *walletUseCase) DeleteOne(ctx context.Context, condition entity.Wallet) error {
	if err := w.repo.DeleteOne(ctx, condition); err != nil {
		log.Print(err)
		return entity.ErrDeleteWallet
	}
	return nil
}
