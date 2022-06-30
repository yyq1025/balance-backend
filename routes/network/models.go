package network

type Network struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Symbol string `json:"symbol"`
}

type Info struct {
	Network
	BlockNumber uint64  `json:"blockNumber"`
	GasPrice    float64 `json:"gasPrice(GWei)"`
}
