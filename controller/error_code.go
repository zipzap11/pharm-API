package controller

import "net/http"

func GetErrorCode(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
