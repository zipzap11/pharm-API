package repository

import (
	"context"
	"errors"

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
	db := getAllProductQuery(r.db.Preload("Category").Model(&model.Product{}), sortFilter)
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

func (r *productRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logrus.Info("find by id repo = ", id)
	var result model.Product
	err := r.db.Debug().Preload("Category").First(&result, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logrus.WithField("id", id).Error(err)
		return nil, err
	}
	return &result, nil
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"product": product,
		}).Error(err)
		return err
	}

	go r.reindexSearchDoc(ctx, int64(product.ID))
	return nil
}

func (r *productRepositoryImpl) UpdateStock(ctx context.Context, productID int64, stock int) error {
	if err := r.db.Debug().Model(&model.Product{}).Where("id = ?", productID).
		Update("stock", stock).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":       ctx,
			"productID": productID,
			"stock":     stock,
		}).Error(err)
	}
	return nil
}

func (r *productRepositoryImpl) DeleteByID(ctx context.Context, id int64) error {
	if err := r.db.Delete(&model.Product{ID: uint(id)}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
			"id":  id,
		}).Error(err)
		return err
	}
	return nil
}

func (r *productRepositoryImpl) reindexSearchDoc(ctx context.Context, id int64) error {
	if err := r.db.Exec(`UPDATE products SET search_text = (
		setweight(to_tsvector(coalesce(name, '')), 'A') || 
		setweight(to_tsvector(coalesce(description, '')), 'B')
	);`).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
			"id":  id,
		}).Error(err)
		return err
	}
	return nil
}
