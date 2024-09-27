package request

type OrderRequest struct {
	TableId   int        `json:"tableId" binding:"required"`
	MenuItems []MenuItem `json:"menuItems" binding:"required"`
}
type MenuItem struct {
	MenuItemID int `json:"menuItemId" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,min=1"`
}
