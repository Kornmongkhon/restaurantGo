package model

type Menus struct {
	MenuItemsId int     `json:"menuItemsId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
