package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	ProvinceID    int64
	StateID       int64
	PostalCode    string
	AddressDetail string
}

type Province struct {
	ID   int64
	Name string
}
type State struct {
	ID   int64
	Name string
}
