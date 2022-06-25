package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
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
		usecase.ErrInvalidEmail,
		usecase.ErrItemAlreadyExist,
		usecase.ErrValidation,
		usecase.ErrAddressNameAlreadyExist,
		usecase.ErrInvalidAddress:
		return http.StatusBadRequest
	case usecase.ErrPermissionDenied:
		return http.StatusForbidden
	default:
		_, ok := err.(validator.ValidationErrors)
		if ok {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}
