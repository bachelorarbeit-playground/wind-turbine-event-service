package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHandler(e *echo.Echo) {
	e.GET("/livez", Livez)
	e.GET("/readyz", Readyz)

	e.GET("/AverageProduction", AverageProduction)

	e.GET("/AnomalyDetection", AnomalyDetection)

}

func Livez(c echo.Context) error {
	return c.String(http.StatusOK, "Server is alive")
}

func Readyz(c echo.Context) error {
	return c.String(http.StatusOK, "Server is ready")
}
