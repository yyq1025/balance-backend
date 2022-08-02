package entity

type Pagination struct {
	IDLte    int `json:"idLte" form:"idLte"`
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}
