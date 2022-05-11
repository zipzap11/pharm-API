package model

import (
	"context"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID int64
	Items  []*CartItem `gorm:"-:all"`
}

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID int64) (*Cart, error)
}

type CartUsecase interface {
	FindCartByUserID(ctx context.Context, userID int64) (*Cart, error)
}
