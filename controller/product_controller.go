package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/zipzap11/pharm-API/dto/request"
	"github.com/zipzap11/pharm-API/util"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/model"
)

type ProductController struct {
	productUsecase model.ProductUsecase
}

func NewProductController(productUsecase model.ProductUsecase) *ProductController {
	return &ProductController{
		productUsecase: productUsecase,
	}
}

func (pc *ProductController) GetAllProducts(c echo.Context) error {
	categoryID, err := getInt64FromQuery("category", c)
	if err != nil {
		logrus.Error(err)
		return c.JSON(http.StatusBadRequest, resp.ErrResponse{
			Message: "invalid category",
		})
	}
	sort := c.QueryParam("sort")
	sortType, err := model.ParseProductSortType(sort)
	if err != nil {
		return c.JSON(http.StatusBadRequest, resp.ErrResponse{
			Message: "invalid sort type",
		})
	}
	query := c.QueryParam("query")
	result, err := pc.productUsecase.GetAllProducts(c.Request().Context(), &model.SortFilter{
		CategoryID: int64(categoryID),
		SortType:   sortType,
		Query:      query,
	})
	if err != nil {
		logrus.Error(err)
		return c.JSON(GetErrorCode(err), resp.ErrResponse{
			Message: err.Error(),
		})
	}
	c.Response().Header().Set("X-Total-Count", fmt.Sprintf("%d", len(result)))
	return c.JSON(http.StatusOK, resp.StdResponse{
		Message: "ok",
		Data:    result,
	})
}

func (pc *ProductController) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	log := logrus.WithField("id", id)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, resp.ErrResponse{
			Message: "invalid id",
		})
	}

	result, err := pc.productUsecase.FindByID(c.Request().Context(), int64(id))
	if err != nil {
		log.Error(err)
		return c.JSON(GetErrorCode(err), resp.ErrResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp.StdResponse{
		Message: "ok",
		Data:    result,
	})
}

func (pc *ProductController) Create(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = &model.Product{}
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	if model.Role(payload.Role) == model.RoleUser {
		return ErrResponseWithCode(c, errors.New("permission denied"), http.StatusForbidden)
	}

	if err := c.Bind(body); err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("invalid payload"), http.StatusBadRequest)
	}

	err := pc.productUsecase.CreateProduct(ctx, body)
	if err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (pc *ProductController) UpdateStock(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = &request.UpdateProductStockRequest{}
	)

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	if model.Role(payload.Role) == model.RoleUser {
		return ErrResponseWithCode(c, errors.New("permission denied"), http.StatusForbidden)
	}

	if err := c.Bind(body); err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("invalid payload"), http.StatusBadRequest)
	}

	v := validator.New()
	fmt.Println()
	if err := v.Struct(body); err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("bad input request"), http.StatusBadRequest)
	}

	if err := pc.productUsecase.UpdateProductStock(ctx, body.ProductID, body.Stock); err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}

func (pc *ProductController) DeleteProduct(c echo.Context) error {
	var (
		ctx     = c.Request().Context()
		log     = logrus.WithField("ctx", ctx)
		idParam = c.Param("id")
		id, err = strconv.ParseInt(idParam, 10, 64)
	)
	if err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, errors.New("bad input request"), http.StatusBadRequest)
	}

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	if model.Role(payload.Role) == model.RoleUser {
		return ErrResponseWithCode(c, errors.New("permission denied"), http.StatusForbidden)
	}

	if err := pc.productUsecase.DeleteProduct(ctx, id); err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}
