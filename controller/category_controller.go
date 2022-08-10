package controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/dto/request"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type categoryController struct {
	categoryUsecase model.CategoryUsecase
}

func NewCategoryController(categoryUsecase model.CategoryUsecase) *categoryController {
	return &categoryController{
		categoryUsecase: categoryUsecase,
	}
}

func (cc *categoryController) GetAllCategories(c echo.Context) error {
	result, err := cc.categoryUsecase.GetAllCategories(c.Request().Context())
	if err != nil {
		return c.JSON(GetErrorCode(err), resp.ErrResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp.StdResponse{
		Message: "ok",
		Data:    resp.ModelToCategoryResponseArray(result),
	})
}

func (cc *categoryController) CreateCategory(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = logrus.WithField("ctx", ctx)
		body = &request.CreateCategoryRequest{}
	)
	if err := c.Bind(body); err != nil {
		log.Error(err)
		return ErrResponseWithCode(c, err, http.StatusBadRequest)
	}

	payload, ok := c.Get(model.AuthorizationPayloadKey).(*util.Payload)
	if !ok {
		return ErrResponseWithCode(c, errors.New("internal server error"), http.StatusInternalServerError)
	}

	if model.Role(payload.Role) == model.RoleUser {
		return ErrResponseWithCode(c, errors.New("permission denied"), http.StatusForbidden)
	}

	if err := cc.categoryUsecase.Create(ctx, body.ToModel()); err != nil {
		log.Error(err)
		return ErrResponse(c, err)
	}

	return SuccessResponse(c, nil)
}
