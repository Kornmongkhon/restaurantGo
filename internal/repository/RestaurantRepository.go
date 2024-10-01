package repository

import (
	"Restaurant/config"
	"Restaurant/database"
	"Restaurant/internal/model"
	"Restaurant/internal/request"
	"database/sql"
	"fmt"
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
	UpdateOrder(tableId int, orderId int, status string) error
	UpdateOrderWithTx(tableId int, orderId int, status string, tx *sql.Tx) error
	DeleteOrder(r *request.OrderRequest) error
	PayOrder(r *request.OrderRequest, tx *sql.Tx) error
	CheckOrderStatus(r *request.OrderRequest, tx *sql.Tx) (string, error)
	HasOrderBeenReviewed(r *request.OrderRequest, tx *sql.Tx) (bool, error)
	ReviewOrder(r *request.OrderRequest, tx *sql.Tx) error
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

func (r *MySQLRestaurantRepository) UpdateOrder(tableId int, orderId int, status string) error {
	var currentStatus string
	currentTime := config.FormatTime(time.Now())
	updateQuery := `
			UPDATE orders
			SET status = ?, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
	_, err := database.DB.Exec(updateQuery, status, currentTime, orderId, tableId)
	if err != nil {
		return err
	}
	statusQuery := "SELECT status FROM orders WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE"
	err = database.DB.QueryRow(statusQuery, orderId, tableId).Scan(&currentStatus)
	log.Println("currentStatus", currentStatus)
	if err != nil {
		return err
	}

	if currentStatus == "canceled" {
		deleteQuery := `
			UPDATE orders
			SET status = ?, is_deleted = TRUE, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
		_, err := database.DB.Exec(deleteQuery, status, currentTime, orderId, tableId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLRestaurantRepository) UpdateOrderWithTx(tableId int, orderId int, status string, tx *sql.Tx) error {
	var currentStatus string
	currentTime := config.FormatTime(time.Now())
	updateQuery := `
			UPDATE orders
			SET status = ?, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
	_, err := tx.Exec(updateQuery, status, currentTime, orderId, tableId)
	if err != nil {
		return err
	}
	statusQuery := "SELECT status FROM orders WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE"
	err = tx.QueryRow(statusQuery, orderId, tableId).Scan(&currentStatus)
	log.Println("currentStatus", currentStatus)
	if err != nil {
		return err
	}

	if currentStatus == "canceled" {
		deleteQuery := `
			UPDATE orders
			SET status = ?, is_deleted = TRUE, updated_at = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
		_, err := tx.Exec(deleteQuery, status, currentTime, orderId, tableId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLRestaurantRepository) DeleteOrder(ro *request.OrderRequest) error {
	if ro.Status == "canceled" {
		currentTime := config.FormatTime(time.Now())
		deleteQuery := `
			UPDATE orders
			SET is_deleted = TRUE, updated_at = ?, status = ?
			WHERE order_id = ? AND table_id = ? AND is_deleted = FALSE;
		`
		_, err := database.DB.Exec(deleteQuery, currentTime, ro.Status, ro.OrderId, ro.TableId)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("cannot delete order, status is not 'canceled'")
}

func (r *MySQLRestaurantRepository) CheckOrderStatus(ro *request.OrderRequest, tx *sql.Tx) (string, error) {
	checkStatusQuery := `
		SELECT status FROM orders
		WHERE order_id = ? AND is_deleted = FALSE
	`

	var status string
	err := tx.QueryRow(checkStatusQuery, ro.OrderId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("order not found or already deleted")
		}
		return "", err
	}
	return status, nil
}

func (r *MySQLRestaurantRepository) PayOrder(ro *request.OrderRequest, tx *sql.Tx) error {
	payQuery := `
		INSERT INTO bills (order_id, table_id, total_amount, bill_date)
		SELECT o.order_id, o.table_id, SUM(oi.quantity * oi.price) AS total_amount, ?
		FROM orders o
		INNER JOIN order_items oi ON o.order_id = oi.order_id
		WHERE o.order_id = ? AND o.is_deleted = FALSE
		GROUP BY o.order_id, o.table_id
	`
	currentTime := config.FormatTime(time.Now())
	_, err := tx.Exec(payQuery, currentTime, ro.OrderId)
	if err != nil {
		return fmt.Errorf("failed to create bill: %v", err)
	}

	return nil
}

func (r *MySQLRestaurantRepository) HasOrderBeenReviewed(ro *request.OrderRequest, tx *sql.Tx) (bool, error) {
	reviewQuery := `
		SELECT count(1) FROM reviews
		WHERE order_id = ? AND is_deleted = FALSE
	`
	var count int
	err := tx.QueryRow(reviewQuery, ro.OrderId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *MySQLRestaurantRepository) ReviewOrder(ro *request.OrderRequest, tx *sql.Tx) error {
	reviewQuery := `
		INSERT INTO reviews (order_id, rating, comment, review_date)
		VALUES (?, ?, ?, ?)
	`
	currentTime := config.FormatTime(time.Now())
	_, err := tx.Exec(reviewQuery, ro.OrderId, ro.Rating, ro.Comment, currentTime)
	if err != nil {
		return fmt.Errorf("failed to create review: %v", err)
	}

	return nil
}
