package repository

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type CartItemRepositoryImpl struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) model.CartItemRepository {
	return &CartItemRepositoryImpl{
		db: db,
	}
}

func (r *CartItemRepositoryImpl) GetItemsByCartID(ctx context.Context, cartID int64) ([]*model.CartItem, error) {
	var items []*model.CartItem
	fmt.Println("cartId = ", cartID)
	err := r.db.Model(&model.CartItem{}).Where("cart_id = ?", cartID).Find(&items).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"cartID": cartID,
		})
		return nil, err
	}

	return items, nil
}
