package repository

import (
	"context"

	"github.com/google/uuid"
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

func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, id).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"id":  id,
			"ctx": ctx,
		}).Error(err)
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, c *model.Category) error {
	c.ID = uuid.New()
	if err := r.db.Create(c).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      ctx,
			"category": c,
		}).Error(err)
		return err
	}
	return nil
}
