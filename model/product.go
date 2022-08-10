package model

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	CategoryID  uuid.UUID      `json:"category_id" validate:"required,uuid"`
	Category    *Category      `json:"category"`
	Name        string         `json:"name" validate:"required"`
	Price       int64          `json:"price" validate:"required"`
	Weight      float64        `json:"weight" validate:"required,numeric"`
	Description string         `json:"description" validate:"required,gt=15"`
	Stock       int            `json:"stock"`
	ImageURL    string         `json:"image_url" validate:"required,url"`
}

type SortFilter struct {
	CategoryID int64
	SortType   ProductSortType
	Query      string
}

type ProductUsecase interface {
	GetAllProducts(ctx context.Context, sortFilter *SortFilter) ([]*Product, error)
	FindByID(ctx context.Context, id int64) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProductStock(ctx context.Context, productID int64, quantity int) error
	DeleteProduct(ctx context.Context, productID int64) error
}

type ProductRepository interface {
	GetAllProducts(ctx context.Context, sortFilter *SortFilter) ([]*Product, error)
	FindByID(ctx context.Context, id int64) (*Product, error)
	Create(ctx context.Context, product *Product) error
	UpdateStock(ctx context.Context, productID int64, quantity int) error
	DeleteByID(ctx context.Context, id int64) error
}

type ProductSortType string

const (
	SortProductAsc  ProductSortType = "ASC"
	SortProductDesc ProductSortType = "DESC"
	SortProductNone ProductSortType = "NONE"
)

func ParseProductSortType(sortType string) (ProductSortType, error) {
	switch {
	case sortType == "asc":
		return SortProductAsc, nil
	case sortType == "desc":
		return SortProductDesc, nil
	case sortType == "":
		return SortProductNone, nil
	default:
		return "", errors.New("invalid sort type")
	}
}
