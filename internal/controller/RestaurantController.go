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
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/ [get]
func (rc *RestaurantController) Home(c echo.Context) error {
	log.Println("RestController -> Healthy Check")
	return service.HealthyCheck(c)
}

// @Summary Find Table
// @Description Find a table by its ID
// @Tags restaurant
// @Accept json
// @Produce json
// @Param table body request.TableRequest true "Table Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/table [post]
func (rc *RestaurantController) FindTable(c echo.Context) error {
	log.Println("RestController -> FindTable")
	var tableRequest request.TableRequest
	if err := c.Bind(&tableRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", tableRequest.TableId)
	responses, status := rc.RestaurantService.FindTable(&tableRequest)
	return c.JSON(status, responses)
}

// @Summary Get all menu
// @Description Retrieve a list of all menu items
// @Tags restaurant
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
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

// @Summary Get order details by table and order ID
// @Description Get detailed information about an order, including menu items
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param order body request.OrderRequest true "Order Request"
// @Success 200 {object} response.CustomResponse
// @Failure 400 {object} response.CustomResponse
// @Failure 500 {object} response.CustomResponse
// @Router /api/v1/restaurant/order/details [post]
func (rc *RestaurantController) OrderDetails(c echo.Context) error {
	log.Println("RestController -> OrderDetails")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	log.Println("OrderID :", orderRequest.OrderId)
	responses, status := rc.RestaurantService.OrderDetails(&orderRequest)
	return c.JSON(status, responses)
}

func (rc *RestaurantController) OrderHistory(c echo.Context) error {
	log.Println("RestController -> OrderHistory")
	var orderRequest request.OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", orderRequest.TableId)
	responses, status := rc.RestaurantService.OrderHistory(&orderRequest)
	return c.JSON(status, responses)
}

func (rc *RestaurantController) UpdateTable(c echo.Context) error {
	log.Println("RestController -> UpdateTable")
	var tableRequest request.TableRequest
	if err := c.Bind(&tableRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", tableRequest.TableId)
	log.Println("Status :", tableRequest.TableStatus)
	responses, status := rc.RestaurantService.UpdateTable(&tableRequest)
	return c.JSON(status, responses)
}

func (rc *RestaurantController) DeleteAllOrderWhenCheckOut(c echo.Context) error {
	log.Println("RestController -> DeleteAllOrderWhenCheckOut")
	var tableRequest request.TableRequest
	if err := c.Bind(&tableRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage(),
		})
	}
	log.Println("TableID :", tableRequest.TableId)
	responses, status := rc.RestaurantService.DeleteAllOrderWhenCheckOut(&tableRequest)
	return c.JSON(status, responses)
}
