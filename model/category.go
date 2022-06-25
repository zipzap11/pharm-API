package model

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string
	ImageURL string
}

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]*Category, error)
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]*Category, error)
}
