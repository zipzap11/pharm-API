package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	ctrl "github.com/zipzap11/pharm-API/controller"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

const (
	authorizationTypeBearer = "bearer"
)

var (
	ErrNoToken      = errors.New("authorization not provided")
	ErrInvalidToken = errors.New("token is invalid")
)

func AuthPaseto(tokenProvider util.TokenProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header["Authorization"]
			if len(authorization) == 0 {
				return ctrl.ErrResponseWithCode(c, ErrNoToken, http.StatusUnauthorized)

			}
			fields := strings.Fields(authorization[0])
			if len(fields) < 2 {
				return ctrl.ErrResponseWithCode(c, ErrInvalidToken, http.StatusUnauthorized)
			}

			authorizationType := fields[0]
			if authorizationType != authorizationTypeBearer {
				return ctrl.ErrResponseWithCode(c, ErrInvalidToken, http.StatusUnauthorized)
			}

			token := fields[1]
			payload, err := tokenProvider.VerifyToken(token)
			if err != nil {
				return ctrl.ErrResponseWithCode(c, err, http.StatusUnauthorized)
			}
			fmt.Println("payload =", payload.ExpiredAt)
			c.Set(model.AuthorizationPayloadKey, payload)

			return next(c)
		}
	}
}
