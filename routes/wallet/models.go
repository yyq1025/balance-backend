package wallet

import (
	"github.com/ethereum/go-ethereum/common"
)

type Wallet struct {
	ID      int            `json:"id"`
	UserID  string         `json:"-"`
	Address common.Address `json:"address"`
	Network string         `json:"network"`
	Token   common.Address `json:"token"`
}

type Balance struct {
	Symbol  string  `json:"symbol"`
	Balance float64 `json:"balance"`
}

type Result struct {
	Wallet
	Balance
}

type Pagination struct {
	IDLte    int `json:"idLte"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
