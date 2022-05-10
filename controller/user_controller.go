package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	resp "github.com/zipzap11/pharm-API/dto/response"

	"github.com/zipzap11/pharm-API/model"
)

type UserController struct {
	userUsecase model.UserUsecase
}

func NewUserController(userUsecase model.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (uc *UserController) CreateUser(c echo.Context) error {
	log := logrus.WithField("ctx", c.Request().Context())
	var body request.CreateUserRequest
	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	err = uc.userUsecase.CreateUser(c.Request().Context(), request.ModelFromCreateUserRequest(&body))
	if err != nil {
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (uc *UserController) Login(c echo.Context) error {
	log := logrus.WithField("ctx", c.Request().Context())

	var body request.LoginRequest
	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	accessToken, refreshToken, err := uc.userUsecase.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, resp.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
