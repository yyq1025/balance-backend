package network

type Network struct {
	ChainID  string `json:"chainId"`
	Name     string `gorm:"primaryKey" json:"name"`
	URL      string `json:"url"`
	Symbol   string `json:"symbol"`
	Explorer string `json:"explorer"`
}
