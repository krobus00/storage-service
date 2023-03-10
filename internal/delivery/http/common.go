package http

import (
	"context"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/labstack/echo/v4"
)

func buildContext(eCtx echo.Context) context.Context {
	token := eCtx.Get(string(constant.KeyTokenCtx))
	ctx := context.WithValue(eCtx.Request().Context(), constant.KeyTokenCtx, token)
	return ctx
}
