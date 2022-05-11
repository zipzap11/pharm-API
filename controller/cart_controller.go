package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type CartController struct {
	cartUsecase model.CartUsecase
}

func NewCartController(cartUsecase model.CartUsecase) *CartController {
	return &CartController{
		cartUsecase: cartUsecase,
	}
}

func (cc *CartController) FindCart(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)
	fmt.Println("payload =", c.Get(model.AuthorizationPayloadKey))
	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}
	fmt.Println("pay =", payload)
	cart, err := cc.cartUsecase.FindCartByUserID(ctx, payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, cart)
}
