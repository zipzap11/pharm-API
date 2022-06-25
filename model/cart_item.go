package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
	Quantity  int64        `json:"quantity"`
	ProductID int64        `json:"product_id"`
	CartID    int64        `json:"cart_id"`
	Product   *Product     `json:"product"`
}

type CartItemRepository interface {
	GetItemsByCartID(ctx context.Context, cartID int64) (items []*CartItem, err error)
	CreateItem(ctx context.Context, cartID, productID int64) error
	FindItemIDByCartIDAndProductID(ctx context.Context, cartID, productID int64) (ID int64, err error)
	FindItemByCartIDAndProductID(ctx context.Context, cartID, productID int64) (item *CartItem, err error)
	UpdateItemQuantity(ctx context.Context, itemID int64, quantity int64) error
	UpdateItemQuantityByType(ctx context.Context, itemID int64, operationType QuantityUpdateType) error
	DeleteItemByID(ctx context.Context, itemID int64) error
	FindItemByID(ctx context.Context, itemID int64) (*CartItem, error)
	DeleteByCartID(ctx context.Context, tx *gorm.DB, cartID int64) error
}

type QuantityUpdateType string

const (
	ADD_QUANTITY_OPERATION_TYPE      = QuantityUpdateType("ADD")
	SUBTRACT_QUANTITY_OPERATION_TYPE = QuantityUpdateType("SUBTRACT")
	UPDATE_QUANTITY_OPERATION_TYPE   = QuantityUpdateType("UPDATE")
)

func ParseQuantityUpdateTypeFromString(in string) (QuantityUpdateType, error) {
	switch in {
	case "ADD":
		return ADD_QUANTITY_OPERATION_TYPE, nil
	case "SUBTRACT":
		return SUBTRACT_QUANTITY_OPERATION_TYPE, nil
	case "UPDATE":
		return UPDATE_QUANTITY_OPERATION_TYPE, nil
	default:
		return "", errors.New("enum type not available")
	}
}
