package http

import (
	"context"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/labstack/echo/v4"
)

func buildContext(eCtx echo.Context) context.Context {
	userID := eCtx.Get(string(constant.KeyUserIDCtx))
	ctx := context.WithValue(eCtx.Request().Context(), constant.KeyUserIDCtx, userID)
	return ctx
}
