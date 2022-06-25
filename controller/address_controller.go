package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type addressController struct {
	addressUsecase model.AddressUsecase
}

func NewAddressController(addressUsecase model.AddressUsecase) *addressController {
	return &addressController{
		addressUsecase: addressUsecase,
	}
}

func (ac *addressController) GetProvinces(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)
	res, err := ac.addressUsecase.GetProvinces(ctx)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}

func (ac *addressController) GetStates(c echo.Context) error {
	var (
		ctx        = c.Request().Context()
		log        = logrus.WithField("ctx", ctx)
		provinceID = c.QueryParam("province")
	)

	if len(provinceID) == 0 {
		return ErrResponseWithCode(c, errors.New("province id is required"), http.StatusBadRequest)
	}
	
	res, err := ac.addressUsecase.GetStatesByProvinceID(ctx, provinceID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}

func (ac *addressController) CreateAddress(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = request.CreateAddressRequest{}
	)

	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	address := request.ModelfromCreateAddressRequest(&body)
	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	address.UserID = payload.UserID

	err = ac.addressUsecase.CreateAddress(ctx, address)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (ac *addressController) GetAddressesByUserID(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	res, err := ac.addressUsecase.GetAddressesByUserID(ctx, payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}
