package model

import "gorm.io/gorm"

type TransactionStatus string

const (
	TransactionStatusPaid    TransactionStatus = "PAID"
	TransactionStatusExpired TransactionStatus = "EXPIRED"
	TransactionStatusPending TransactionStatus = "PENDING"
)

type Transaction struct {
	gorm.Model
	UserID          int64
	Status          TransactionStatus
	Resi            string
	Price           int64
	PaymentMethodID int64
	PaymentURL      string
	Shipping        string
	AddressID       int64
}
