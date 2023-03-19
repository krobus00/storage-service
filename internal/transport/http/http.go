package http

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type HTTPDelivery struct {
	e                *echo.Echo
	objectController *ObjectController
}

func NewHTTPDelivery() *HTTPDelivery {
	return new(HTTPDelivery)
}

// InjectEcho :nodoc:
func (t *HTTPDelivery) InjectEcho(e *echo.Echo) error {
	if e == nil {
		return errors.New("invalid echo")
	}
	t.e = e
	return nil
}

// InjectUserController :nodoc:
func (t *HTTPDelivery) InjectObjectController(c *ObjectController) error {
	if c == nil {
		return errors.New("invalid object controller")
	}
	t.objectController = c
	return nil
}

func (t *HTTPDelivery) InitRoutes() {
	api := t.e.Group("/api")

	storage := api.Group("/storage")
	storage.GET("/", t.objectController.GetPresignURL, DecodeJWTToken(true))
	storage.POST("/upload", t.objectController.Upload, DecodeJWTToken(false))
}
