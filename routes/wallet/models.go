package wallet

import "github.com/ethereum/go-ethereum/common"

type Wallet struct {
	Id      int            `gorm:"autoIncrement" json:"id"`
	UserId  int            `json:"-"`
	Address common.Address `json:"address"`
	Network string         `gorm:"default:Ethereum" json:"network"`
	Token   common.Address `json:"token"`
	Tag     string         `json:"tag,omitempty"`
}

type Balance struct {
	Wallet
	Symbol  string `json:"symbol"`
	Balance any    `json:"balance"`
}
