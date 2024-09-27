package service

import (
	"Restaurant/database"
	"Restaurant/internal/repository"
	"Restaurant/internal/request"
	"Restaurant/internal/response"
	"Restaurant/utils/enums"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

type RestaurantService struct {
	RestaurantRepo repository.RestaurantRepository
}

func HealthyCheck(c echo.Context) error {
	log.Println("RestaurantService -> Healthy Check")
	return c.JSON(http.StatusOK, response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	})
}

func (s *RestaurantService) GetAllMenu() (response.CustomResponse, int) {
	log.Println("RestaurantService -> GetAllMenu")
	menus, err := s.RestaurantRepo.GetAllMenu()
	if err != nil {
		log.Printf("Service error fetching menus: %v", err)
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage(),
		}, http.StatusNotFound
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
		Data:    menus}, http.StatusOK
}

func (s *RestaurantService) OrderMenu(c *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> OrderMenu")
	//check input
	if c.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if len(c.MenuItems) == 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", menuItems must not be empty.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", menuItems must not be empty.",
		}, http.StatusBadRequest
	}
	//find table id
	existsTableId, err := s.RestaurantRepo.FindTableById(c)
	if err != nil {
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	if !existsTableId {
		log.Println("RestaurantService -> " + enums.NotFound.GetMessage() + ", Table ID not found.")
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage() + ", Table ID not found.",
		}, http.StatusNotFound
	}
	//find menu id
	notFoundItems, err := s.RestaurantRepo.FindMenuItemById(c.MenuItems)
	if err != nil {
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	if len(notFoundItems) > 0 {
		joinIDS := joinWithComma(notFoundItems)
		log.Println("RestaurantService -> "+enums.NotFound.GetMessage()+", MenuItem IDs not found: ", joinIDS)
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage() + ", MenuItem IDs not found: " + joinIDS,
		}, http.StatusNotFound
	}
	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	orderId, err := s.RestaurantRepo.InsertOrder(c, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error inserting order:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	err = s.RestaurantRepo.InsertOrderItems(orderId, c.MenuItems, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error inserting order items:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	err = tx.Commit()
	if err != nil {
		log.Println("RestaurantService -> Error committing transaction:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}
func joinWithComma(ids []int) string {
	notFoundItemsStr := make([]string, len(ids))
	for i, id := range ids {
		notFoundItemsStr[i] = fmt.Sprint(id)
	}
	return strings.Join(notFoundItemsStr, ", ")
}
