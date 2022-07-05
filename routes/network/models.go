package network

type Network struct {
	ChainId  string `json:"chainId"`
	Name     string `gorm:"primaryKey" json:"name"`
	Url      string `json:"url"`
	Symbol   string `json:"symbol"`
	Explorer string `json:"explorer"`
}

type Info struct {
	Network
	BlockNumber uint64  `json:"blockNumber"`
	GasPrice    float64 `json:"gasPrice(GWei)"`
}
