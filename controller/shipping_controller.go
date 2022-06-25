package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type ShippingController struct {
	shippingUsecase model.ShippingUsecase
}

func NewShippingController(shippingUsecase model.ShippingUsecase) *ShippingController {
	return &ShippingController{
		shippingUsecase: shippingUsecase,
	}
}

func (ac *ShippingController) GetShippingPackages(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		log     = logrus.WithField("ctx", ctx)
		address = c.QueryParam("address")
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	addressID, err := strconv.Atoi(address)
	if err != nil {
		return ErrResponseWithCode(c, errors.New("invalid address_id"), http.StatusBadRequest)
	}

	res, err := ac.shippingUsecase.GetShippingPackages(ctx, int64(addressID), payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}
