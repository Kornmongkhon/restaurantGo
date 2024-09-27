package controller

import (
	"Restaurant/internal/request"
	"Restaurant/internal/response"
	"Restaurant/internal/service"
	"Restaurant/utils/enums"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type RestaurantController struct {
	RestaurantService *service.RestaurantService
}

func (rc *RestaurantController) Home(c echo.Context) error {
	log.Println("RestController -> Healthy Check")
	return service.HealthyCheck(c)
}

func (rc *RestaurantController) GetAllMenu(c echo.Context) error {
	log.Println("RestController -> GetAllMenu")
	responses := rc.RestaurantService.GetAllMenu()
	return c.JSON(http.StatusOK, responses)
}

func (rc *RestaurantController) OrderMenu(c echo.Context) error {
	log.Println("RestController -> OrderMenu")
	var orderRequest request.OrderRequest

	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	println("TableID : ", orderRequest.TableId)
	responses := rc.RestaurantService.OrderMenu(&orderRequest)
	return c.JSON(http.StatusOK, responses)
}
