package controller

import (
	"Restaurant/service"
	"github.com/labstack/echo/v4"
	"log"
)

func Home(c echo.Context) error {
	log.Println("RestController -> Healthy Check")
	return service.HealthyCheck(c)
}

func GetAllMenu(c echo.Context) error {
	log.Println("RestController -> GetAllMenu")

	return service.GetAllMenu(c)
}
