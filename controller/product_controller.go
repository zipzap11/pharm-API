package controller

import (
	"net/http"

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

	return c.JSON(http.StatusOK, resp.StdResponse{
		Message: "ok",
		Data:    result,
	})
}
