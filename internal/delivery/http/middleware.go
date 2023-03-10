package http

import (
	"net/http"
	"strings"

	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/labstack/echo/v4"
)

// DecodeJWTToken :nodoc:
func DecodeJWTToken(allowGuest bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			res := model.NewResponse().WithMessage(model.ErrTokenInvalid.Error())
			token := eCtx.Request().Header.Get("Authorization")
			token = strings.Replace(token, "Bearer ", "", -1)
			if token == "" && !allowGuest {
				return eCtx.JSON(http.StatusUnauthorized, res)
			}

			eCtx.Set(string(constant.KeyTokenCtx), token)
			return next(eCtx)
		}
	}
}
