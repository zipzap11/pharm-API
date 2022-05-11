package model

import (
	"context"

	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	Quantity  int64
	ProductID int64
	CartID    int64
}

type CartItemRepository interface {
	GetItemsByCartID(ctx context.Context, cartID int64) (items []*CartItem, err error)
}
