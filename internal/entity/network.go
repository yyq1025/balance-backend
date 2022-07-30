package entity

import "context"

type Network struct {
	ChainID  string `json:"chainId"`
	Name     string `gorm:"primaryKey" json:"name"`
	URL      string `json:"url"`
	Symbol   string `json:"symbol"`
	Explorer string `json:"explorer"`
}

type NetworkUseCase interface {
	GetAll(ctx context.Context) ([]Network, error)
}

type NetworkRepository interface {
	GetAll(ctx context.Context, networks *[]Network) error
}
