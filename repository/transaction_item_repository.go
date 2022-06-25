package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type transactionItemRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionItemRepository(db *gorm.DB) model.TransactionItemRepository {
	return &transactionItemRepositoryImpl{
		db: db,
	}
}

func (r *transactionItemRepositoryImpl) GetItemsByTransactionID(ctx context.Context, transactionID int64) ([]*model.TransactionItem, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":           ctx,
		"transactionID": transactionID,
	})

	var items []*model.TransactionItem
	err := r.db.Where("transaction_id = ?", transactionID).
		Preload("Product").
		Find(&items).
		Error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return items, nil
}

func (r *transactionItemRepositoryImpl) CreateMultipleItems(ctx context.Context, tx *gorm.DB, transactionID int64, items []*model.TransactionItem) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":           ctx,
		"items":         items,
		"transactionID": transactionID,
	})

	err := tx.Create(&items).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
