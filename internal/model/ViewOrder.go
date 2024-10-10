package model

type ViewOrder struct {
	OrderId   int    `json:"orderId"`
	TableId   int    `json:"tableId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}
