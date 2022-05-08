package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) model.CategoryRepository {
	return &CategoryRepositoryImpl{
		db: db,
	}
}

func (r *CategoryRepositoryImpl) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	var result []*model.Category
	err := r.db.Model(&model.Category{}).Find(&result).Error
	if err != nil {
		logrus.WithField("ctx", ctx).Error(err)
		return nil, err
	}
	return result, nil
}
