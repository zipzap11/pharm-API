package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type cartUsecaseImpl struct {
	cartRepository model.CartRepository
}

func NewCartUsecase(cartRepository model.CartRepository) model.CartUsecase {
	return &cartUsecaseImpl{
		cartRepository: cartRepository,
	}
}

func (u *cartUsecaseImpl) FindCartByUserID(ctx context.Context, userID int64) (*model.Cart, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})
	cart, err := u.cartRepository.FindCartByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if cart == nil {
		return nil, ErrNotFound
	}

	return cart, nil
}
