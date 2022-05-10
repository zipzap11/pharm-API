package model

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Phone    string
}

type UserUsecase interface {
	CreateUser(ctx context.Context, user *User) error
	Login(ctx context.Context, email string, password string) (string, string, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}
