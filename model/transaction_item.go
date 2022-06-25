package model

import (
	"context"

	"gorm.io/gorm"
)

type TransactionItem struct {
	TransactionID int64    `json:"transaction_id"`
	ProductID     int64    `json:"product_id"`
	Quantity      int64    `json:"quantity"`
	Product       *Product `json:"product"`
}

type TransactionItemRepository interface {
	GetItemsByTransactionID(ctx context.Context, transactionID int64) ([]*TransactionItem, error)
	CreateMultipleItems(ctx context.Context, tx *gorm.DB, transactionID int64, items []*TransactionItem) error
}
