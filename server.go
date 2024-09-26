package main

import (
	"Restaurant/database"
	"github.com/labstack/echo/v4"

	"Restaurant/controller"
)

func main() {
	e := echo.New()
	database.InitDB("user:password@tcp(127.0.0.1:3306)/restaurant?parseTime=true")
	apiV1 := e.Group("/api/v1/restaurant")
	apiV1.GET("/all/menu", controller.GetAllMenu)
	apiV1.GET("/", controller.Home)
	apiV1.POST("/order/menu", controller.OrderMenu)
	e.Logger.Fatal(e.Start(":1323"))
}
