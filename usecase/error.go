package usecase

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidEmail       = errors.New("email is invalid")
	ErrInvalidCredential  = errors.New("wrong email or password")
	ErrIncorrectUserToken = errors.New("wrong user")
	ErrBlockedSession     = errors.New("session is blocked")
	ErrMissmatchedToken   = errors.New("missmatched token")
)
