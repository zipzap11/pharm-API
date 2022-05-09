package usecase

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidEmail      = errors.New("email is invalid")
	ErrInvalidCredential = errors.New("wrong email or password")
)
