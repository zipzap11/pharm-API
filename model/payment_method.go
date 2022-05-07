package model

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	Name string
}
