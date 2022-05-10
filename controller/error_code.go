package controller

import (
	"net/http"

	"github.com/zipzap11/pharm-API/usecase"
)

func GetErrorCode(err error) int {
	switch err {
	case usecase.ErrNotFound:
		return http.StatusNotFound
	case usecase.ErrBlockedSession,
		usecase.ErrIncorrectUserToken,
		usecase.ErrMissmatchedToken:
		return http.StatusUnauthorized
	case usecase.ErrInvalidCredential,
		usecase.ErrInvalidEmail:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
