package request

type OrderRequest struct {
	OrderId   int        `json:"orderId" binding:"required"`
	TableId   int        `json:"tableId" binding:"required"`
	Status    string     `json:"status" binding:"required"`
	MenuItems []MenuItem `json:"menuItems" binding:"required"`
	Rating    int        `json:"rating" binding:"required,max=5 min=1"`
	Comment   string     `json:"comment" binding:"required"`
}
type MenuItem struct {
	MenuItemID int `json:"menuItemId" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,min=1"`
}
