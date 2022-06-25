package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name          string   `json:"name" validate:"required"`
	UserID        int64    `json:"user_id"`
	ProvinceID    string   `json:"province_id" validate:"required"`
	StateID       string   `json:"state_id" validate:"required"`
	Province      Province `json:"province"`
	State         State    `json:"state"`
	PostalCode    string   `json:"postal_code" validate:"min=6,max=6"`
	AddressDetail string   `json:"address_detail" validate:"required"`
}

type Province struct {
	ID   string `json:"province_id"`
	Name string `json:"province"`
}

type State struct {
	ID         string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"city_name"`
	Type       string `json:"type"`
}

type AddressRepository interface {
	GetProvinces(ctx context.Context) ([]*Province, error)
	GetStatesByProvinceID(ctx context.Context, provinceID string) ([]*State, error)
	GetProvinceByIDFromAPI(ctx context.Context, provinceID string) (*Province, error)
	GetStateByIDFromAPI(ctx context.Context, stateID string) (*State, error)
	GetProvinceByIDFromDB(ctx context.Context, provinceID string) (*Province, error)
	GetStateByIDFromDB(ctx context.Context, stateID string) (*State, error)
	CreateAddress(ctx context.Context, address *Address) error
	CreateState(ctx context.Context, state *State) error
	CreateProvince(ctx context.Context, province *Province) error
	GetAddressesByUserID(ctx context.Context, userID int64) ([]*Address, error)
	GetAddressByID(ctx context.Context, id int64) (*Address, error)
	GetAddressByNameAndUserID(ctx context.Context, name string, userID int64) (*Address, error)
}

type AddressUsecase interface {
	GetProvinces(ctx context.Context) ([]*Province, error)
	GetStatesByProvinceID(ctx context.Context, provinceID string) ([]*State, error)
	CreateAddress(ctx context.Context, address *Address) error
	GetAddressesByUserID(ctx context.Context, userID int64) ([]*Address, error)
	GetAddressByID(ctx context.Context, id int64) (*Address, error)
}
