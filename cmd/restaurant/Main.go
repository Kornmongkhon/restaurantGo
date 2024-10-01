package main

import (
	"Restaurant/config"
	"Restaurant/database"
	_ "Restaurant/docs"
	"Restaurant/internal/controller"
	"Restaurant/internal/repository"
	"Restaurant/internal/response"
	"Restaurant/internal/service"
	"Restaurant/utils/enums"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"log"
	"net/http"
)

func main() {
	cfg := config.DBLoadConfig()
	config.SetTimeZone("Asia/Bangkok")
	dataSourceName := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true"
	database.InitDB(dataSourceName)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c) // เรียก handler ถัดไป
			if err != nil {
				// พิมพ์ข้อผิดพลาดที่เกิดขึ้น
				log.Printf("Error occurred in %s %s: %v", c.Request().Method, c.Path(), err)

				// ตรวจสอบประเภทของข้อผิดพลาดและสร้าง CustomResponse
				var customResponse response.CustomResponse
				switch err.(type) {
				case *echo.HTTPError:
					httpError := err.(*echo.HTTPError)
					customResponse.Code = enums.Error.GetCode()
					customResponse.Message = enums.Error.GetMessage()
					return c.JSON(httpError.Code, customResponse)
				default:
					customResponse.Code = enums.Error.GetCode()
					customResponse.Message = enums.Error.GetMessage()
					return c.JSON(http.StatusInternalServerError, customResponse)
				}
			}
			return nil
		}
	})
	restaurantRepo := &repository.MySQLRestaurantRepository{}
	restaurantService := &service.RestaurantService{RestaurantRepo: restaurantRepo}
	restaurantController := &controller.RestaurantController{RestaurantService: restaurantService}
	apiV1 := e.Group("/api/v1/restaurant")
	apiV1.GET("/swagger/*", echoSwagger.WrapHandler)
	apiV1.GET("/", restaurantController.Home)
	apiV1.GET("/all/menu", restaurantController.GetAllMenu)
	apiV1.POST("/order/menu", restaurantController.OrderMenu)
	apiV1.PATCH("/order/update", restaurantController.UpdateOrder)
	apiV1.DELETE("/order/delete", restaurantController.DeleteOrder)
	apiV1.POST("/order/pay", restaurantController.PayOrder)
	apiV1.POST("/order/review", restaurantController.ReviewOrder)
	e.Logger.Fatal(e.Start(":1323"))
}
