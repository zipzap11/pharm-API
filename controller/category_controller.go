package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	resp "github.com/zipzap11/pharm-API/dto/response"
	"github.com/zipzap11/pharm-API/model"
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
		Data:    result,
	})
}
