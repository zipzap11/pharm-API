package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/db"
	"github.com/zipzap11/pharm-API/model"
)

type ProductController struct{}

func (c *ProductController) GetAllProducts(ctx echo.Context) error {
	var result []model.Product
	err := db.DB.Model(&model.Product{}).Find(&result).Error
	if err != nil {
		logrus.Error(err)
		return ctx.String(http.StatusInternalServerError, "internal server error")
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    result,
	})
}
