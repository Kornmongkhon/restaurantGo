package model

type Menus struct {
	MenuItemsId int          `json:"menuItemsId"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	IsAvailable bool         `json:"isAvailable"`
	FileObjects []FileObject `json:"fileObjects"`
}

type FileObject struct {
	FileName string `json:"fileName"`
	Base64   string `json:"base64"`
}
