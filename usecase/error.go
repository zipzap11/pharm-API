package usecase

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrInvalidCredential       = errors.New("wrong email or password")
	ErrIncorrectUserToken      = errors.New("wrong user")
	ErrBlockedSession          = errors.New("session is blocked")
	ErrMissmatchedToken        = errors.New("missmatched token")
	ErrItemAlreadyExist        = errors.New("item already exist")
	ErrPermissionDenied        = errors.New("operation not allowed")
	ErrValidation              = errors.New("bad input request")
	ErrCreateTransaction       = errors.New("failed to create transaction")
	ErrAddressNameAlreadyExist = errors.New("address name already exist")
	ErrInvalidAddress          = errors.New("address is invalid")
)
