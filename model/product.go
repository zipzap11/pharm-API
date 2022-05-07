package model

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	CategoryID  int64
	Name        string
	Price       int64
	Weight      int64
	Description string
}

type SortFilter struct {
	CategoryID int64
	SortType   ProductSortType
	Query      string
}

type ProductUsecase interface {
	GetAllProducts(ctx context.Context, sortFilter *SortFilter) ([]*Product, error)
}

type ProductRepository interface {
	GetAllProducts(ctx context.Context, sortFilter *SortFilter) ([]*Product, error)
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
