package request

type OrderRequest struct {
	TableId   int        `json:"tableId" binding:"required"`
	MenuItems []MenuItem `json:"menu_items" binding:"required"`
}
type MenuItem struct {
	MenuItemID int     `json:"menuItemId" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required"`
}
