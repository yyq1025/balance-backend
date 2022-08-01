package entity

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type Wallet struct {
	ID          int            `json:"id"`
	UserID      string         `json:"-"`
	Address     common.Address `json:"address"`
	NetworkName string         `json:"networkName"`
	Token       common.Address `json:"token"`
	Network     Network        `gorm:"foreignKey:NetworkName" json:"network"`
}

type Balance struct {
	Wallet
	Symbol  string  `json:"symbol"`
	Balance float64 `json:"balance"`
}

type WalletUseCase interface {
	AddOne(ctx context.Context, wallet *Wallet) (Balance, error)
	GetOne(ctx context.Context, userID string, id int) (Balance, error)
	GetManyWithPagination(ctx context.Context, userID string, pagination *Pagination) ([]Balance, *Pagination, error)
	DeleteOne(ctx context.Context, userID string, ID int) error
}

type WalletRepository interface {
	AddOne(ctx context.Context, wallet *Wallet) error
	GetOne(ctx context.Context, userID string, id int, wallet *Wallet) error
	GetManyWithPagination(ctx context.Context, userID string, pagination *Pagination, wallets *[]Wallet) error
	DeleteOne(ctx context.Context, userID string, id int) error
}

type WalletEthAPI interface {
	GetSymbol(ctx context.Context, wallet Wallet) (string, error)
	GetBalance(ctx context.Context, wallet Wallet) (float64, error)
}
