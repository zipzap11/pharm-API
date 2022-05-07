package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryID  int64
	Name        string
	Price       int64
	Weight      int64
	Description string
}
