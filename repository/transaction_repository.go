package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type transactionRepositoryImpl struct {
	db                        *gorm.DB
	transactionItemRepository model.TransactionItemRepository
}

func NewTransactionRepository(db *gorm.DB, transactionItemRepository model.TransactionItemRepository) model.TransactionRepository {
	return &transactionRepositoryImpl{
		db:                        db,
		transactionItemRepository: transactionItemRepository,
	}
}

func (r *transactionRepositoryImpl) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) (int64, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":         ctx,
		"transaction": transaction,
	})

	err := tx.Create(transaction).Error
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return int64(transaction.ID), nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatus(ctx context.Context, transactionID int64, status model.TransactionStatus) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":           ctx,
		"transactionID": transactionID,
		"status":        status,
	})

	err := r.db.Model(&model.Transaction{}).
		Where("id = ?", transactionID).
		Update("status", status).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *transactionRepositoryImpl) GetTransactionByUserID(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	var (
		transactions []*model.Transaction
	)
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})

	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error
	if err != nil {
		log.Error(err)
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	ch := make(chan *model.Transaction)
	for _, v := range transactions {
		g.Go(func() error {
			v := <-ch
			items, err := r.transactionItemRepository.GetItemsByTransactionID(ctx, int64(v.ID))
			if err != nil {
				log.Error(err)
				return err
			}
			v.Items = items
			return nil
		})
		ch <- v
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepositoryImpl) UpdateTransactionPaymentURL(ctx context.Context, tx *gorm.DB, transactionID int64, paymentURL string) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":           ctx,
		"transactionID": transactionID,
		"paymentURL":    paymentURL,
	})
	log.Info("url =", paymentURL)
	err := tx.Model(&model.Transaction{}).
		Where("id = ?", transactionID).
		Update("payment_url", paymentURL).
		Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
