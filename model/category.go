package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`
	Name      string       `json:"name"`
	ImageURL  string       `json:"image_url"`
}

type CategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]*Category, error)
	Create(ctx context.Context, c *Category) error
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]*Category, error)
	FindByID(ct context.Context, id int64) (*Category, error)
	Create(ctx context.Context, c *Category) error
}
