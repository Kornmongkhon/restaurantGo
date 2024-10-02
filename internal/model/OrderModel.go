package model

type Order struct {
	OrderId    int          `json:"orderId"`
	TableId    int          `json:"tableId"`
	Status     string       `json:"status"`
	OrderItems []OrderItems `json:"orderItems"`
}

type OrderItems struct {
	MenuItemId  int     `json:"menuItemId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
