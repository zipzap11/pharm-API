package repository

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DEFAULT_QUANTITY = 1

type cartItemRepositoryImpl struct {
	db *gorm.DB
}

func NewCartItemRepository(db *gorm.DB) model.CartItemRepository {
	return &cartItemRepositoryImpl{
		db: db,
	}
}

func (r *cartItemRepositoryImpl) GetItemsByCartID(ctx context.Context, cartID int64) ([]*model.CartItem, error) {
	var items []*model.CartItem
	err := r.db.
		Preload("Product").
		Where("cart_id = ?", cartID).
		Find(&items).
		Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"cartID": cartID,
		})
		return nil, err
	}
	return items, nil
}

func (r *cartItemRepositoryImpl) CreateItem(ctx context.Context, cartID, productID int64) error {
	err := r.db.Create(&model.CartItem{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  DEFAULT_QUANTITY,
	}).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"cartID": cartID,
		})
		return err
	}
	return nil
}

func (r *cartItemRepositoryImpl) FindItemIDByCartIDAndProductID(ctx context.Context, cartID, productID int64) (int64, error) {
	var item model.CartItem
	err := r.db.Model(&model.CartItem{}).
		Select("id").
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).
		Error
	if err == nil {
		return int64(item.ID), nil
	}

	if errors.Is(gorm.ErrRecordNotFound, err) {
		return 0, nil
	}
	logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"cartID":    cartID,
		"productID": productID,
	})
	return 0, err
}

func (r *cartItemRepositoryImpl) UpdateItemQuantity(ctx context.Context, itemID int64, quantity int64) error {
	err := r.db.Model(&model.CartItem{}).
		Where("id = ?", itemID).
		Update("quantity", quantity).
		Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      ctx,
			"itemID":   itemID,
			"quantity": quantity,
		}).Error(err)
		return err
	}
	return nil
}

func (r *cartItemRepositoryImpl) UpdateItemQuantityByType(ctx context.Context, itemID int64, operationType model.QuantityUpdateType) error {
	var updateClause clause.Expr
	switch operationType {
	case model.ADD_QUANTITY_OPERATION_TYPE:
		updateClause = gorm.Expr("quantity + 1")
	case model.SUBTRACT_QUANTITY_OPERATION_TYPE:
		updateClause = gorm.Expr("quantity - 1")
	}
	err := r.db.Model(&model.CartItem{}).
		Where("id = ?", itemID).
		Update("quantity", updateClause).
		Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"itemID":  itemID,
			"opsType": operationType,
		}).Error(err)
		return err
	}
	return nil
}

func (r *cartItemRepositoryImpl) FindItemByCartIDAndProductID(ctx context.Context, cartID, productID int64) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("cartID = ? AND productID = ?", cartID, productID).
		First(&item).
		Error
	if err == nil {
		return &item, nil
	}

	switch {
	case errors.Is(gorm.ErrRecordNotFound, err):
		return nil, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx":       ctx,
			"cartID":    cartID,
			"productID": productID,
		}).Error(err)
		return nil, err
	}
}

func (r *cartItemRepositoryImpl) DeleteItemByID(ctx context.Context, itemID int64) error {
	err := r.db.Delete(&model.CartItem{}, itemID).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"itemID": itemID,
		}).Error(err)
		return err
	}
	return nil
}

func (r *cartItemRepositoryImpl) FindItemByID(ctx context.Context, itemID int64) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.First(&item, itemID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"itemID": itemID,
		}).Error(err)
		return nil, err
	}
	return &item, nil
}

func (r *cartItemRepositoryImpl) DeleteByCartID(ctx context.Context, tx *gorm.DB, cartID int64) error {
	err := tx.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    ctx,
			"cartID": cartID,
		}).Error(err)
		return err
	}
	return nil
}
