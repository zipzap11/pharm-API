package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
	Name      string       `validate:"required" json:"name"`
	Email     string       `validate:"required,email" json:"email"`
	Password  string       `validate:"required,min=6" json:"password"`
	Phone     string       `json:"phone"`
}

type UserUsecase interface {
	CreateUser(ctx context.Context, user *User) error
	Login(ctx context.Context, email string, password string) (string, string, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, tx *gorm.DB, user *User) (id int64, err error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, userID int64) (*User, error)
}
