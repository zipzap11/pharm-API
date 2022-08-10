package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type productUsecaseImpl struct {
	productRepository  model.ProductRepository
	categoryRepository model.CategoryRepository
}

func NewProductUsecase(productRepository model.ProductRepository, categoryRepository model.CategoryRepository) model.ProductUsecase {
	return &productUsecaseImpl{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

func (u *productUsecaseImpl) GetAllProducts(ctx context.Context, sortFilter *model.SortFilter) ([]*model.Product, error) {
	result, err := u.productRepository.GetAllProducts(ctx, sortFilter)
	if err != nil {
		logrus.WithField("sortFilter", sortFilter).Error(err)
		return nil, err
	}
	return result, err
}

func (u *productUsecaseImpl) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	logrus.Info("find by id usecase = ", id)
	result, err := u.productRepository.FindByID(ctx, id)
	if err != nil {
		logrus.WithField("id", id).Error(err)
		return nil, err
	}
	if result == nil {
		return nil, ErrNotFound
	}
	return result, nil
}

func (u *productUsecaseImpl) CreateProduct(ctx context.Context, product *model.Product) error {
	err := u.productRepository.Create(ctx, product)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"product": product,
		}).Error(err)
		return err
	}
	return nil
}

func (u *productUsecaseImpl) UpdateProductStock(ctx context.Context, productID int64, stock int) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"productID": productID,
		"stock":     stock,
	})
	_, err := u.FindByID(ctx, productID)
	if err != nil {
		log.Error(err)
		return err
	}
	err = u.productRepository.UpdateStock(ctx, productID, stock)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":   ctx,
			"stock": stock,
		}).Error(err)
		return err
	}
	return nil
}

func (u *productUsecaseImpl) DeleteProduct(ctx context.Context, productID int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"productID": productID,
	})

	_, err := u.FindByID(ctx, productID)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := u.productRepository.DeleteByID(ctx, productID); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
