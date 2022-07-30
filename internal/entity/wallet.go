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
	*Wallet
	Symbol  string  `json:"symbol"`
	Balance float64 `json:"balance"`
}

type WalletUseCase interface {
	AddOne(context.Context, *Wallet) (Balance, error)
	GetOne(context.Context, *Wallet) (Balance, error)
	GetManyWithPagination(context.Context, *Wallet, *Pagination) ([]Balance, *Pagination, error)
	DeleteOne(context.Context, *Wallet) error
}

type WalletRepository interface {
	AddOne(ctx context.Context, wallet *Wallet) error
	GetOne(ctx context.Context, condition, wallet *Wallet) error
	GetManyWithPagination(ctx context.Context, condition *Wallet, wallet *[]Wallet, p *Pagination) error
	DeleteOne(ctx context.Context, condition *Wallet) error
}
