package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/util"

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
		// return ErrResponse(c, err)
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
	}

	return SuccessResponse(c, nil)
}

func (uc *UserController) CreateSuperUser(c echo.Context) error {
	log := logrus.WithField("ctx", c.Request().Context())
	var body request.CreateUserRequest
	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	err = uc.userUsecase.CreateSuperUser(c.Request().Context(), request.ModelFromCreateUserRequest(&body))
	if err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, err.Error())
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

func (uc *UserController) FindCurrentUser(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	log.Info("ok = ", ok)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}
	log.Info("payload = ", payload)
	res, err := uc.userUsecase.FindByID(ctx, payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}

func (uc *UserController) FindUsers(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	if payload.Role == int(model.RoleUser) {
		return ErrResponseWithCode(c, errors.New("permission denied"), http.StatusForbidden)
	}

	res, err := uc.userUsecase.FindAllUsers(ctx)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, resp.UserResponsesFromModel(res))
}