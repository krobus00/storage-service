package http

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type Delivery struct {
	e                *echo.Echo
	objectController *ObjectController
}

func NewDelivery() *Delivery {
	return new(Delivery)
}

func (t *Delivery) InjectEcho(e *echo.Echo) error {
	if e == nil {
		return errors.New("invalid echo")
	}
	t.e = e
	return nil
}

func (t *Delivery) InjectObjectController(c *ObjectController) error {
	if c == nil {
		return errors.New("invalid object controller")
	}
	t.objectController = c
	return nil
}

func (t *Delivery) InitRoutes() {
	api := t.e.Group("/api")

	storage := api.Group("/storage")
	storage.GET("/", t.objectController.GetPresignURL, DecodeJWTToken(true))
	storage.POST("/upload", t.objectController.Upload, DecodeJWTToken(false))
}
