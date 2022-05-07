package model

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	Quantity  int64
	ProductID int64
	CartID    int64
}
