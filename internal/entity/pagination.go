package entity

type Pagination struct {
	IDLte    int `json:"idLte"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
