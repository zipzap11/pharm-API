package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type TransactionController struct {
	transactionUsecase model.TransactionUsecase
}

func NewTransactionController(transactionUsecase model.TransactionUsecase) *TransactionController {
	return &TransactionController{
		transactionUsecase: transactionUsecase,
	}
}

func (tc *TransactionController) GetTotalPrice(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		log      = logrus.WithField("ctx", ctx)
		address  = c.QueryParam("address")
		shipping = c.QueryParam("shipping")
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	addressID, err := strconv.Atoi(address)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("address id not valid"), http.StatusBadRequest)
	}

	price, shippingPrice, err := tc.transactionUsecase.GetTotalPrice(ctx, payload.UserID, int64(addressID), shipping)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, resp.TotalPriceResponse{
		Price:         price,
		ShippingPrice: shippingPrice,
		TotalPrice:    price + shippingPrice,
	})
}

func (tc *TransactionController) CreateTransaction(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = request.CreateTransactionRequest{}
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

	res, err := tc.transactionUsecase.CreateTransaction(ctx, payload.UserID, body.AddressID, body.ShippingServices)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, res)
}

func (tc *TransactionController) HandleTransactionCallback(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = map[string]interface{}{}
	)

	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	orderID, exist := body["order_id"].(string)
	if !exist {
		return ErrResponseWithCode(c, errors.New("id not exist"), http.StatusBadRequest)
	}

	transactionID, err := strconv.Atoi(orderID)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("id not valid"), http.StatusBadRequest)
	}

	err = tc.transactionUsecase.UpdateTransactionStatus(ctx, int64(transactionID))
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (tc *TransactionController) GetTransactionByUserID(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		log = logrus.WithField("ctx", ctx)
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	res, err := tc.transactionUsecase.GetTransactionByUserID(ctx, payload.UserID)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, model.ToTransactionResponses(res))
}
