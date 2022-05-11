package repository

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type cartRepositoryImpl struct {
	db                 *gorm.DB
	cartItemRepository model.CartItemRepository
}

func NewCartRepository(db *gorm.DB, cartItemRepository model.CartItemRepository) model.CartRepository {
	return &cartRepositoryImpl{
		db:                 db,
		cartItemRepository: cartItemRepository,
	}
}

func (r *cartRepositoryImpl) FindCartByUserID(ctx context.Context, userID int64) (*model.Cart, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})
	var cart model.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	items, err := r.cartItemRepository.GetItemsByCartID(ctx, int64(cart.ID))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	cart.Items = items

	return &cart, nil
}
