package service

import (
	"Restaurant/database"
	"Restaurant/enums"
	"Restaurant/model/model"
	"Restaurant/model/response"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func HealthyCheck(c echo.Context) error {
	log.Println("MenuService -> Healthy Check")
	return c.JSON(http.StatusOK, response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.String(),
	})
}

func GetAllMenu(c echo.Context) error {
	log.Println("MenuService -> GetAllMenu")

	query := "select id, name, description, price from menu_items where is_available = true"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.String(),
		})
	}
	defer rows.Close()

	var menus []model.Menus

	// Iterate over the result set and scan each row into a Menu struct
	for rows.Next() {
		var menu model.Menus
		if err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price); err != nil {
			log.Printf("Error scanning row: %v", err)
			return c.JSON(http.StatusInternalServerError, response.CustomResponse{
				Code:    enums.Invalid.GetCode(),
				Message: enums.Invalid.String(),
			})
		}
		menus = append(menus, menu)
	}

	// Check for errors after iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.String(),
		})
	}

	// Log successful data retrieval
	log.Println("MenuService -> Menus successfully retrieved")

	return c.JSON(http.StatusOK, response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.String(),
		Data:    menus,
	})
}

func OrderMenu(c echo.Context) error {
	log.Println("MenuService -> OrderMenu")
	//query := "insert into orders (table_id, total_amount,create_at ) values (?, ?, ?)"
	//rows, err := database.DB.Query(query, 1, 1000, "2021-09-01")
	//if err != nil {
	//	log.Printf("Error executing query: %v", err)
	//	return c.JSON(http.StatusInternalServerError, response.CustomResponse{
	//		Code:    enums.Error.GetCode(),
	//		Message: enums.Error.String(),
	//	})
	//}
	//defer rows.Close()
	return c.JSON(http.StatusOK, response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.String(),
	})
}
