package repository

import (
	"Restaurant/database"
	"Restaurant/internal/model"
	"log"
)

type RestaurantRepository interface {
	GetAllMenu() ([]model.Menus, []model.Menus)
}
type MySQLRestaurantRepository struct{}

func (r *MySQLRestaurantRepository) GetAllMenu() ([]model.Menus, []model.Menus) {
	query := "SELECT id, name, description, price FROM menu_items WHERE is_available = true"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching menus from database: %v", err)
		return nil, nil
	}
	defer rows.Close()

	var menus []model.Menus
	for rows.Next() {
		var menu model.Menus
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price); err != nil {
			log.Printf("Error scanning menu: %v", err)
			return nil, nil
		}
		menus = append(menus, menu)
	}

	return menus, nil
}
