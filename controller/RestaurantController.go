package controller

import (
	"Restaurant/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the Restaurant API")
}

func GetAllMenu(c echo.Context) error {
	restaurants := service.GetAllMenu()
	return c.JSON(http.StatusOK, restaurants)

}
