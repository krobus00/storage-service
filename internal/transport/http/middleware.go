package http

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/krobus00/storage-service/internal/constant"
	"github.com/krobus00/storage-service/internal/model"
	"github.com/labstack/echo/v4"
)

func DecodeJWTToken(allowGuest bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(eCtx echo.Context) error {
			res := model.NewResponse().WithMessage(model.ErrTokenInvalid.Error())
			accessToken := eCtx.Request().Header.Get("Authorization")
			accessToken = strings.ReplaceAll(accessToken, "Bearer ", "")
			if accessToken == "" && !allowGuest {
				return eCtx.JSON(http.StatusUnauthorized, res)
			}

			token, _ := jwt.Parse(accessToken, nil)

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return eCtx.JSON(http.StatusUnauthorized, res)
			}
			userID, ok := claims["userID"]
			if !ok {
				return eCtx.JSON(http.StatusUnauthorized, res)
			}

			eCtx.Set(string(constant.KeyUserIDCtx), userID)
			return next(eCtx)
		}
	}
}
