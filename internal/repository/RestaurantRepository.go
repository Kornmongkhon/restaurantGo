package repository

import (
	"Restaurant/config"
	"Restaurant/database"
	"Restaurant/internal/model"
	"Restaurant/internal/request"
	"database/sql"
	"log"
	"time"
)

type RestaurantRepository interface {
	GetAllMenu() ([]model.Menus, error)
	FindTableById(c *request.OrderRequest) (bool, error)
	FindMenuItemById(c []request.MenuItem) ([]int, error)
	InsertOrder(c *request.OrderRequest, tx *sql.Tx) (int64, error)
	InsertOrderItems(orderID int64, menuItems []request.MenuItem, tx *sql.Tx) error
	FindOrderById(r *request.OrderRequest) (bool, error)
	UpdateOrder(r *request.OrderRequest) error
}
type MySQLRestaurantRepository struct{}

func (r *MySQLRestaurantRepository) GetAllMenu() ([]model.Menus, error) {
	query := "SELECT menu_items_id, name, description, price FROM menu_items WHERE is_available = true"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching menus from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var menus []model.Menus
	for rows.Next() {
		var menu model.Menus
		if err := rows.Scan(&menu.MenuItemsId, &menu.Name, &menu.Description, &menu.Price); err != nil {
			log.Printf("Error scanning menu: %v", err)
			return nil, err
		}
		menus = append(menus, menu)
	}

	return menus, nil
}

func (r *MySQLRestaurantRepository) FindTableById(c *request.OrderRequest) (bool, error) {
	query := "SELECT count(1) FROM tables WHERE table_id = ? AND is_deleted = FALSE"
	var count int
	err := database.DB.QueryRow(query, c.TableId).Scan(&count)
	if count > 0 {
		return true, nil
	}

	return false, err
}

func (r *MySQLRestaurantRepository) FindMenuItemById(c []request.MenuItem) ([]int, error) {
	query := "SELECT count(1) FROM menu_items WHERE menu_items_id = ?  AND is_deleted = FALSE AND is_available = TRUE;"
	var notFoundItems []int
	for _, item := range c {
		var count int
		err := database.DB.QueryRow(query, item.MenuItemID).Scan(&count)
		if err != nil {
			log.Printf("Error checking menu item ID %d: %v", item.MenuItemID, err)
			return nil, err
		}
		if count == 0 {
			notFoundItems = append(notFoundItems, item.MenuItemID)
		}
	}

	return notFoundItems, nil
}

func (r *MySQLRestaurantRepository) InsertOrder(c *request.OrderRequest, tx *sql.Tx) (int64, error) {
	orderQuery := "INSERT INTO orders (table_id, created_at) VALUES (?, ?)"
	currentTime := config.FormatTime(time.Now())
	result, err := tx.Exec(orderQuery, c.TableId, currentTime)
	if err != nil {
		return 0, err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *MySQLRestaurantRepository) InsertOrderItems(orderID int64, menuItems []request.MenuItem, tx *sql.Tx) error {
	for _, menuItem := range menuItems {
		menuItemQuery := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, price)
		SELECT ?, ?, ?, price FROM menu_items WHERE menu_items_id = ? AND is_available = TRUE AND is_deleted = FALSE
	`
		_, err := tx.Exec(menuItemQuery, orderID, menuItem.MenuItemID, menuItem.Quantity, menuItem.MenuItemID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLRestaurantRepository) FindOrderById(ro *request.OrderRequest) (bool, error) {
	query := "SELECT count(1) FROM orders WHERE order_id = ? AND is_deleted = FALSE AND status NOT IN ('canceled')"
	var count int
	err := database.DB.QueryRow(query, ro.OrderId).Scan(&count)
	if count > 0 {
		return true, nil
	}

	return false, err
}

func (r *MySQLRestaurantRepository) UpdateOrder(ro *request.OrderRequest) error {
	var currentStatus string
	currentTime := config.FormatTime(time.Now())
	updateQuery := `
			UPDATE orders
			SET status = ?, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
	_, err := database.DB.Exec(updateQuery, ro.Status, currentTime, ro.OrderId, ro.TableId)
	if err != nil {
		return err
	}
	statusQuery := "SELECT status FROM orders WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE"
	err = database.DB.QueryRow(statusQuery, ro.OrderId, ro.TableId).Scan(&currentStatus)
	log.Println("currentStatus", currentStatus)
	if err != nil {
		return err
	}

	if currentStatus == "canceled" {
		deleteQuery := `
			UPDATE orders
			SET is_deleted = TRUE, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
		_, err := database.DB.Exec(deleteQuery, currentTime, ro.OrderId, ro.TableId)
		if err != nil {
			return err
		}
	}

	return nil
}
