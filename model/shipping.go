package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Shipping struct {
	ID         uint         `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  sql.NullTime `gorm:"index" json:"deleted_at"`
	AddressID   int64 `json:"address_id"`
	Services    string `json:"services"`
	Description string `json:"description"`
	ETD         string `json:"etd"`
	Resi        string `json:"resi"`
	Price       int64 `json:"price"`
}

type ROCostResponse struct {
	RO struct {
		Results []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Costs []struct {
				Service     string `json:"service"`
				Description string `json:"description"`
				Cost        []struct {
					Value int64  `json:"value"`
					ETD   string `json:"etd"`
					Note  string `json:"note"`
				} `json:"cost"`
			} `json:"costs"`
		}
	} `json:"rajaOngkir"`
}

type ShippingRepository interface {
	GetShippingsPackages(ctx context.Context, stateID string, weight float64) ([]*Shipping, error)
	CreateShipping(ctx context.Context, tx *gorm.DB, shipping *Shipping) (id int64, err error)
}

type ShippingUsecase interface {
	GetShippingPackages(ctx context.Context, addressID, userID int64) ([]*Shipping, error)
	GetShippingPackageByServices(ctx context.Context, addressID, userID int64, services string) (*Shipping, error)
	CreateShipping(ctx context.Context, tx *gorm.DB, shipping *Shipping) (id int64, err error)
}
