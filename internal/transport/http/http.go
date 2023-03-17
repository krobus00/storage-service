package http

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type HTTPDelivery struct {
	e                 *echo.Echo
	storageController *StorageController
}

func NewHTTPDelivery() *HTTPDelivery {
	return new(HTTPDelivery)
}

// InjectEcho :nodoc:
func (d *HTTPDelivery) InjectEcho(e *echo.Echo) error {
	if e == nil {
		return errors.New("invalid echo")
	}
	d.e = e
	return nil
}

// InjectUserController :nodoc:
func (d *HTTPDelivery) InjectStorageController(c *StorageController) error {
	if c == nil {
		return errors.New("invalid storage controller")
	}
	d.storageController = c
	return nil
}

func (d *HTTPDelivery) InitRoutes() {
	api := d.e.Group("/api")

	storage := api.Group("/storage")
	storage.GET("/", d.storageController.GetPresignURL, DecodeJWTToken(true))
	storage.POST("/upload", d.storageController.Upload, DecodeJWTToken(false))
}
