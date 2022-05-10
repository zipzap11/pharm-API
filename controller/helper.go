package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	resp "github.com/zipzap11/pharm-API/dto/response"
)

func getInt64FromQuery(key string, c echo.Context) (int64, error) {
	val := c.QueryParam(key)
	if len(val) == 0 {
		return 0, nil
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return int64(res), nil
}

func ErrResponse(c echo.Context, err error) error {
	return c.JSON(GetErrorCode(err), resp.ErrResponse{
		Message: err.Error(),
	})
}

func ErrResponseWithCode(c echo.Context, err error, code int) error {
	return c.JSON(code, resp.ErrResponse{
		Message: err.Error(),
	})
}

func SuccessResponse(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, resp.StdResponse{
		Message: "ok",
		Data:    data,
	})
}
