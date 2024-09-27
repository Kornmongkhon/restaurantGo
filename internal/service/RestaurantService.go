package service

import (
	"Restaurant/internal/repository"
	"Restaurant/internal/request"
	"Restaurant/internal/response"
	"Restaurant/utils/enums"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
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

func (s *RestaurantService) GetAllMenu() response.CustomResponse {
	log.Println("RestaurantService -> GetAllMenu")
	menus, err := s.RestaurantRepo.GetAllMenu()
	if err != nil {
		log.Printf("Service error fetching menus: %v", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
		Data:    menus,
	}
}

func (s *RestaurantService) OrderMenu(c *request.OrderRequest) response.CustomResponse {
	log.Println("RestaurantService -> OrderMenu")
	// Implement the logic to order menu
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}
}
