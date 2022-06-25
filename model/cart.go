package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID         uint         `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  sql.NullTime `gorm:"index" json:"deleted_at"`
	UserID     int64        `json:"user_id"`
	TotalPrice int64        `json:"total_price"`
	Items      []*CartItem  `gorm:"-:all" json:"items"`
}

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID int64) (*Cart, error)
	GetCartIDByUserID(ctx context.Context, userID int64) (cartID int64, err error)
	CreateCart(ctx context.Context, tx *gorm.DB, userID int64) error
}

type CartUsecase interface {
	FindCartByUserID(ctx context.Context, userID int64) (*Cart, error)
	AddItemToCart(ctx context.Context, userID, productID int64) error
	UpdateItemQuantity(ctx context.Context, itemID, userID, quantity int64, opsType QuantityUpdateType) error
	RemoveItemFromCart(ctx context.Context, itemID, userID int64) error
	GetCartTotalPrice(ctx context.Context, userID int64) (price int64, err error)
}
