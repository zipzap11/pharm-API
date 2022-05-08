package controller

import (
	"net/http"

	"github.com/zipzap11/pharm-API/usecase"
)

func GetErrorCode(err error) int {
	switch err {
	case usecase.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
