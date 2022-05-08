package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type categoryUsecaseImpl struct {
	categoryRepository model.CategoryRepository
}

func NewCategoryUsecase(categoryRepository model.CategoryRepository) model.CategoryUsecase {
	return &categoryUsecaseImpl{
		categoryRepository: categoryRepository,
	}
}

func (u *categoryUsecaseImpl) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	result, err := u.categoryRepository.GetAllCategories(ctx)
	if err != nil {
		logrus.WithField("ctx", ctx).Error(err)
		return nil, err
	}
	return result, nil
}
