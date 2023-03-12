package handler

import (
	"net/http"

	"event-sourcing/pkg/storage"
	"github.com/labstack/echo/v4"
)

func AverageProduction(c echo.Context) error {
	data := storage.GetMaterializedViewEntries(storage.Store.AverageProduction)
	return c.JSON(http.StatusOK, data)
}

func AnomalyDetection(c echo.Context) error {
	data := storage.GetMaterializedViewEntries(storage.Store.AnomalyDetection)
	return c.JSON(http.StatusOK, data)
}
