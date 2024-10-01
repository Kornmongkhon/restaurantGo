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

// @Summary Health Check
// @Description Check the health status of the server
// @Success 200 {string} string "OK"
// @Router /api/v1/restaurant/ [get]
func (rc *RestaurantController) Home(c echo.Context) error {
	log.Println("RestController -> Healthy Check")
	return service.HealthyCheck(c)
}

// @Summary Get all menu
// @Description Retrieve a list of all menu items
// @Tags restaurant
// @Success 200 {object} response.CustomResponse
// @Router /api/v1/restaurant/all/menu [get]
func (rc *RestaurantController) GetAllMenu(c echo.Context) error {
	log.Println("RestController -> GetAllMenu")
	responses, status := rc.RestaurantService.GetAllMenu()
	return c.JSON(status, responses)
}

// @Summary Order Menu
// @Description Place an order for menu items
// @Tags restaurant
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/menu [post]
func (rc *RestaurantController) OrderMenu(c echo.Context) error {
	log.Println("RestController -> OrderMenu")
	var orderRequest request.OrderRequest

	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	for _, menuItem := range orderRequest.MenuItems {
		log.Println("MenuItemID :", menuItem.MenuItemID, "Quantity :", menuItem.Quantity)
	}
	responses, status := rc.RestaurantService.OrderMenu(&orderRequest)
	return c.JSON(status, responses)
}

// @Summary Update order
// @Description Update the order status
// @Tags restaurant
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/update [patch]
func (rc *RestaurantController) UpdateOrder(c echo.Context) error {
	log.Println("RestController -> UpdateOrder")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	log.Println("OrderID :", orderRequest.OrderId)
	log.Println("Status :", orderRequest.Status)
	responses, status := rc.RestaurantService.UpdateOrder(&orderRequest)
	return c.JSON(status, responses)
}

// @Summary Delete order
// @Description Delete the order by its ID and Table ID
// @Tags restaurant
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/delete [delete]
func (rc *RestaurantController) DeleteOrder(c echo.Context) error {
	log.Println("RestController -> DeleteOrder")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	log.Println("OrderID :", orderRequest.OrderId)
	responses, status := rc.RestaurantService.DeleteOrder(&orderRequest)
	return c.JSON(status, responses)
}

// @Summary Pay for order
// @Description Process payment for the order by its ID and Table ID
// @Tags restaurant
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/pay [post]
func (rc *RestaurantController) PayOrder(c echo.Context) error {
	log.Println("RestController -> PayOrder")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	log.Println("OrderID :", orderRequest.OrderId)
	responses, status := rc.RestaurantService.PayOrder(&orderRequest)
	return c.JSON(status, responses)
}

// @Summary Submit a review for an order
// @Description Add a rating and comment for the specified order
// @Tags restaurant
// @Accept json
// @Produce json
// @Param orderRequest body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/review [post]
func (rc *RestaurantController) ReviewOrder(c echo.Context) error {
	log.Println("RestController -> ReviewOrder")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("OrderID :", orderRequest.OrderId)
	log.Println("Rating :", orderRequest.Rating)
	log.Println("Comment :", orderRequest.Comment)
	responses, status := rc.RestaurantService.ReviewOrder(&orderRequest)
	return c.JSON(status, responses)
}
