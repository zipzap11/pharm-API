package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) GetAllProducts(ctx context.Context, sortFilter *model.SortFilter) ([]*model.Product, error) {
	var (
		result []*model.Product
		err    error
	)
	db := getAllProductQuery(r.db.Model(&model.Product{}), sortFilter)
	err = db.Find(&result).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        ctx,
			"sortFilter": sortFilter,
		})
		return nil, err
	}
	return result, nil
}
