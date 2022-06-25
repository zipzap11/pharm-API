package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
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
	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	cart, err := cc.cartUsecase.FindCartByUserID(ctx, payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, cart)
}

func (cc *CartController) AddItemToCart(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		log     = logrus.WithField("ctx", ctx)
		product = c.QueryParam("product")
	)
	productID, err := strconv.Atoi(product)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("bad product id"), http.StatusBadRequest)
	}

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	err = cc.cartUsecase.AddItemToCart(ctx, payload.UserID, int64(productID))
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (cc *CartController) RemoveItemFromCart(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)

	item := c.QueryParam("item")
	itemID, err := strconv.Atoi(item)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	err = cc.cartUsecase.RemoveItemFromCart(ctx, int64(itemID), payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (cc *CartController) UpdateItemQuantity(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = request.UpdateItemQuantityRequest{}
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	opsType, err := model.ParseQuantityUpdateTypeFromString(body.Type)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	err = cc.cartUsecase.UpdateItemQuantity(ctx, body.ItemID, payload.UserID, body.Quantity, opsType)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}
