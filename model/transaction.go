package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TransactionStatusSuccess TransactionStatus = "SUCCESS"
	TransactionStatusExpired TransactionStatus = "EXPIRED"
	TransactionStatusPending TransactionStatus = "PENDING"
)

type Transaction struct {
	ID              uint               `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	DeletedAt       sql.NullTime       `gorm:"index" json:"deleted_at"`
	UserID          int64              `json:"user_id"`
	User            *User              `json:"user"`
	Status          TransactionStatus  `json:"status"`
	Price           int64              `json:"price"`
	PaymentMethodID int64              `json:"payment_method_id"`
	PaymentMethod   PaymentMethod      `json:"payment_method"`
	PaymentURL      string             `json:"payment_url" gorm:"column:payment_url"`
	ShippingID      int64              `json:"shipping_id"`
	Shipping        *Shipping          `json:"shipping"`
	Items           []*TransactionItem `json:"items"`
}

type TransactionResponse struct {
	Transaction
	CreatedAtStr string `json:"created_at_str"`
}

func (t *Transaction) ToTransactionResponse() *TransactionResponse {
	year, month, day := t.CreatedAt.Date()
	return &TransactionResponse{
		Transaction:  *t,
		CreatedAtStr: fmt.Sprintf("%d %s %02d", year, month.String(), day),
	}
}

func ToTransactionResponses(transactions []*Transaction) []*TransactionResponse {
	var res []*TransactionResponse
	for _, v := range transactions {
		res = append(res, v.ToTransactionResponse())
	}
	return res
}

type TransactionUsecase interface {
	GetTotalPrice(ctx context.Context, userID int64, addressID int64, shippingPackages string) (price, shippingPrice int64, err error)
	CreateTransaction(ctx context.Context, userID, addressID int64, shippingPackages string) (string, error)
	UpdateTransactionStatus(ctx context.Context, transactionID int64) error
	GetTransactionByUserID(ctx context.Context, userID int64) ([]*Transaction, error)
	FindAllTransactions(ctx context.Context) ([]*Transaction, error)
	FindByID(ctx context.Context, id int64, userID int64) (*Transaction, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *Transaction) (id int64, err error)
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status TransactionStatus) error
	GetTransactionByUserID(ctx context.Context, userID int64) ([]*Transaction, error)
	UpdateTransactionPaymentURL(ctx context.Context, tx *gorm.DB, transactionID int64, paymentURL string) error
	FindAll(ctx context.Context) ([]*Transaction, error)
	FindByID(ctx context.Context, id int64) (*Transaction, error)
}

func ParseTransactionStatusFromString(s string) TransactionStatus {
	switch s {
	case "capture", "settlement":
		return TransactionStatusSuccess
	case "cancel", "expire":
		return TransactionStatusExpired
	default:
		return TransactionStatusPending
	}
}
