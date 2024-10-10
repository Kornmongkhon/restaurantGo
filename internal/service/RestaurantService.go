package service

import (
	"Restaurant/database"
	"Restaurant/internal/repository"
	"Restaurant/internal/request"
	"Restaurant/internal/response"
	"Restaurant/utils/enums"
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
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

func (s *RestaurantService) FindTable(r *request.TableRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> FindTable")
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	exists, err := s.RestaurantRepo.FindTableByTableRequestId(r)
	if err != nil {
		log.Printf("Service error fetching table: %v", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	if !exists {
		log.Println("RestaurantService -> " + enums.NotFound.GetMessage() + ", Table ID not found.")
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage() + ", Table ID " + fmt.Sprint(r.TableId) + " not found.",
		}, http.StatusNotFound
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
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
	for i, menu := range menus {
		for j, fileObject := range menu.FileObjects {
			filePath := fmt.Sprintf("K:\\IdeaProjects\\GoLand\\26Sep\\Restaurant\\assets\\images\\%s", fileObject.FileName)
			base64Content, err := s.convertFileToBase64(filePath)
			if err != nil {
				log.Printf("Error converting file to base64: %v", err)
				return response.CustomResponse{
					Code:    enums.Error.GetCode(),
					Message: enums.Error.GetMessage(),
				}, http.StatusInternalServerError
			}
			menus[i].FileObjects[j].Base64 = base64Content
		}
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
	// Check quantity of menu items
	for _, menuItem := range c.MenuItems {
		if menuItem.Quantity <= 0 {
			log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", MenuItem ID " + fmt.Sprint(menuItem.MenuItemID) + ", quantity must be greater than 0.")
			return response.CustomResponse{
				Code:    enums.Invalid.GetCode(),
				Message: enums.Invalid.GetMessage() + ", MenuItem ID " + fmt.Sprint(menuItem.MenuItemID) + ", quantity must be greater than 0.",
			}, http.StatusBadRequest
		}
	}
	//find table id
	resp, status, err := s.CheckTableId(c)
	if err != nil {
		return resp, status
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
	log.Println("Transaction committed successfully")
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}

func (s *RestaurantService) UpdateOrder(r *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> UpdateOrder")
	//check input
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.OrderId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Order ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.Status == "" {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Status must not be empty.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Status must not be empty.",
		}, http.StatusBadRequest
	}
	//find table id
	resp, status, err := s.CheckTableId(r)
	if err != nil {
		return resp, status
	}
	//find order id
	respOrder, status, err := s.CheckOrderId(r)
	if err != nil {
		return respOrder, status
	}
	err = s.RestaurantRepo.UpdateOrder(r.TableId, r.OrderId, r.Status)
	if err != nil {
		log.Println("RestaurantService -> Error updating order:", err)
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}

func (s *RestaurantService) DeleteOrder(r *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> DeleteOrder")
	//check input
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.OrderId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Order ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	//find table id
	resp, status, err := s.CheckTableId(r)
	if err != nil {
		return resp, status
	}
	//find order id
	respOrder, status, err := s.CheckOrderId(r)
	if err != nil {
		return respOrder, status
	}
	err = s.RestaurantRepo.DeleteOrder(r)
	if err != nil {
		log.Println("RestaurantService -> Error deleting order:", err)
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}

func (s *RestaurantService) PayOrder(r *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> PayOrder")
	//check input
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.OrderId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Order ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	//find table id
	resp, status, err := s.CheckTableId(r)
	if err != nil {
		return resp, status
	}
	//find order id
	respOrder, status, err := s.CheckOrderId(r)
	if err != nil {
		return respOrder, status
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
	// Check order status
	statusOrder, err := s.RestaurantRepo.CheckOrderStatus(r, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error checking order status:", err)
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: "Order not found or already deleted.",
		}, http.StatusNotFound
	}

	// Check if status is not "completed"
	if statusOrder != "completed" {
		tx.Rollback()
		log.Println("RestaurantService -> Order is not in 'completed' status")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order is not completed. Cannot proceed with payment.",
		}, http.StatusBadRequest
	}
	err = s.RestaurantRepo.PayOrder(r, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error paying order:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	err = s.RestaurantRepo.UpdateOrderWithTx(r.TableId, r.OrderId, "paid", tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error updating order:", err)
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
	log.Println("Transaction committed successfully")
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}

func (s *RestaurantService) ReviewOrder(r *request.OrderRequest) (response.CustomResponse, int) {
	//check input
	if r.OrderId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Order ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.Rating < 1 || r.Rating > 5 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Rating must be between 1 and 5.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Rating must be between 1 and 5.",
		}, http.StatusBadRequest
	}
	//find order id
	respOrder, status, err := s.CheckOrderId(r)
	if err != nil {
		return respOrder, status
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
	// Check order status
	statusOrder, err := s.RestaurantRepo.CheckOrderStatus(r, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error checking order status:", err)
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: "Order not found or already deleted.",
		}, http.StatusNotFound
	}

	// Check if status is not "completed"
	if statusOrder != "paid" {
		tx.Rollback()
		log.Println("RestaurantService -> Order is not in 'paid' status")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order is not paid. Cannot proceed with payment.",
		}, http.StatusBadRequest
	}
	// Check if the order has already been reviewed
	hasReviewed, err := s.RestaurantRepo.HasOrderBeenReviewed(r, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error checking if order has been reviewed:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	if hasReviewed {
		tx.Rollback()
		log.Println("RestaurantService -> Order has already been reviewed")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order has already been reviewed.",
		}, http.StatusBadRequest
	}
	err = s.RestaurantRepo.ReviewOrder(r, tx)
	if err != nil {
		tx.Rollback()
		log.Println("RestaurantService -> Error reviewing order:", err)
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
	log.Println("Transaction committed successfully")
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
	}, http.StatusOK
}

func (s *RestaurantService) OrderDetails(r *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> OrderDetails")
	//check input
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	if r.OrderId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Order ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Order ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	//find table id
	resp, status, err := s.CheckTableId(r)
	if err != nil {
		return resp, status
	}
	//find order id
	respOrder, status, err := s.CheckOrderId(r)
	if err != nil {
		return respOrder, status
	}
	orderDetails, err := s.RestaurantRepo.GetOrderDetails(r)
	if err != nil {
		log.Println("RestaurantService -> Error getting order details:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
		Data:    orderDetails,
	}, http.StatusOK
}

func (s *RestaurantService) OrderHistory(r *request.OrderRequest) (response.CustomResponse, int) {
	log.Println("RestaurantService -> OrderHistory")
	//check input
	if r.TableId <= 0 {
		log.Println("RestaurantService -> " + enums.Invalid.GetMessage() + ", Table ID must be greater than 0.")
		return response.CustomResponse{
			Code:    enums.Invalid.GetCode(),
			Message: enums.Invalid.GetMessage() + ", Table ID must be greater than 0.",
		}, http.StatusBadRequest
	}
	//find table id
	resp, status, err := s.CheckTableId(r)
	if err != nil {
		return resp, status
	}
	orders, err := s.RestaurantRepo.GetOrderHistory(r)
	if err != nil {
		log.Println("RestaurantService -> Error getting order history:", err)
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError
	}
	return response.CustomResponse{
		Code:    enums.Success.GetCode(),
		Message: enums.Success.GetMessage(),
		Data:    orders,
	}, http.StatusOK
}

func joinWithComma(ids []int) string {
	notFoundItemsStr := make([]string, len(ids))
	for i, id := range ids {
		notFoundItemsStr[i] = fmt.Sprint(id)
	}
	return strings.Join(notFoundItemsStr, ", ")
}

func (s *RestaurantService) CheckTableId(r *request.OrderRequest) (response.CustomResponse, int, error) {
	existsTableId, err := s.RestaurantRepo.FindTableById(r)
	if err != nil {
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError, err
	}
	if !existsTableId {
		log.Println("RestaurantService -> " + enums.NotFound.GetMessage() + ", Table ID not found.")
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage() + ", Table ID " + fmt.Sprint(r.TableId) + " not found.",
		}, http.StatusNotFound, fmt.Errorf("table id not found")
	}
	return response.CustomResponse{}, http.StatusOK, nil
}

func (s *RestaurantService) CheckOrderId(r *request.OrderRequest) (response.CustomResponse, int, error) {
	existsOrderId, err := s.RestaurantRepo.FindOrderById(r)
	if err != nil {
		return response.CustomResponse{
			Code:    enums.Error.GetCode(),
			Message: enums.Error.GetMessage(),
		}, http.StatusInternalServerError, err
	}
	if !existsOrderId {
		log.Println("RestaurantService -> " + enums.NotFound.GetMessage() + ", Order ID not found.")
		return response.CustomResponse{
			Code:    enums.NotFound.GetCode(),
			Message: enums.NotFound.GetMessage() + ", Order ID " + fmt.Sprint(r.OrderId) + " not found.",
		}, http.StatusNotFound, fmt.Errorf("order id not found")
	}
	return response.CustomResponse{}, http.StatusOK, nil
}
func (s *RestaurantService) convertFileToBase64(filePath string) (string, error) {
	fmt.Printf("Attempting to read file at: %s\n", filePath) // Debug print
	data, err := os.ReadFile(filePath)                       // Read the file
	if err != nil {
		return "", err
	}
	base64Content := base64.StdEncoding.EncodeToString(data)
	mimeType := http.DetectContentType(data)
	return "data:" + mimeType + ";base64," + base64Content, nil
}
