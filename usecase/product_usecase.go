package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type productUsecaseImpl struct {
	productRepository model.ProductRepository
}

func NewProductUsecase(productRepository model.ProductRepository) model.ProductUsecase {
	return &productUsecaseImpl{
		productRepository: productRepository,
	}
}

func (u *productUsecaseImpl) GetAllProducts(ctx context.Context, sortFilter *model.SortFilter) ([]*model.Product, error) {
	result, err := u.productRepository.GetAllProducts(ctx, sortFilter)
	if err != nil {
		logrus.WithField("sortFilter", sortFilter)
		return nil, err
	}
	return result, err
}
