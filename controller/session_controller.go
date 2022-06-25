package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/model"
)

type SessionController struct {
	sessionUsecase model.SessionUsecase
}

func NewSessionController(sessionUsecase model.SessionUsecase) *SessionController {
	return &SessionController{
		sessionUsecase: sessionUsecase,
	}
}

func (sc *SessionController) RefreshSession(c echo.Context) error {
	ctx := c.Request().Context()
	log := logrus.WithField("ctx", ctx)

	var body request.RefreshSessionRequest
	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	accessToken, refreshToken, err := sc.sessionUsecase.RefreshSession(ctx, body.RefreshToken)
	if err != nil {
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, resp.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (sc *SessionController) CheckSession(c echo.Context) error {
	ctx := c.Request().Context()
	log := logrus.WithField("ctx", ctx)

	bearer := c.Request().Header["Authorization"]
	fmt.Println("bearer = ", bearer)
	fmt.Println("len =", len(bearer))
	if len(bearer) < 1 {
		return ErrResponseWithCode(c, errors.New("token is required"), http.StatusBadRequest)
	}
	fields := strings.Fields(bearer[0])
	if len(fields) < 2 {
		return ErrResponseWithCode(c, errors.New("invalid token format"), http.StatusBadRequest)
	}
	if fields[0] != "bearer" {
		return ErrResponseWithCode(c, errors.New("invalid token format"), http.StatusBadRequest)
	}

	headerToken := fields[1]
	err := sc.sessionUsecase.CheckSession(ctx, headerToken)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}
