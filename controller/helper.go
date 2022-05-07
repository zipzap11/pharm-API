package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
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
