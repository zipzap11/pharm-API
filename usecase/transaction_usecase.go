package usecase

import (
	"context"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type transactionUsecaseImpl struct {
	cartUsecase               model.CartUsecase
	shippingUsecase           model.ShippingUsecase
	midtransSnapClient        *snap.Client
	midtransCoreClient        *coreapi.Client
	transactionRepository     model.TransactionRepository
	transactionItemRepository model.TransactionItemRepository
	userRepository            model.UserRepository
	addressUsecase            model.AddressUsecase
	db                        *gorm.DB
	cartItemRepository        model.CartItemRepository
}

func NewTransactionUsecase(
	cartUsecase model.CartUsecase,
	shippingUsecase model.ShippingUsecase,
	midtransSnapClient *snap.Client,
	midtransCoreClient *coreapi.Client,
	transactionRepository model.TransactionRepository,
	userRepository model.UserRepository,
	addressUsecase model.AddressUsecase,
	db *gorm.DB,
	transactionItemRepository model.TransactionItemRepository,
	cartItemRepository model.CartItemRepository,
) model.TransactionUsecase {
	return &transactionUsecaseImpl{
		cartUsecase:               cartUsecase,
		shippingUsecase:           shippingUsecase,
		midtransSnapClient:        midtransSnapClient,
		midtransCoreClient:        midtransCoreClient,
		transactionRepository:     transactionRepository,
		userRepository:            userRepository,
		addressUsecase:            addressUsecase,
		db:                        db,
		transactionItemRepository: transactionItemRepository,
		cartItemRepository:        cartItemRepository,
	}
}

func (u *transactionUsecaseImpl) GetTotalPrice(ctx context.Context, userID int64, addressID int64, shippingPackages string) (int64, int64, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":              ctx,
		"userID":           userID,
		"addressID":        addressID,
		"shippingPackages": shippingPackages,
	})

	price, err := u.cartUsecase.GetCartTotalPrice(ctx, userID)
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}

	shipping, err := u.shippingUsecase.GetShippingPackageByServices(ctx, addressID, userID, shippingPackages)
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}

	return price, shipping.Price, nil
}

func (u *transactionUsecaseImpl) CreateTransaction(ctx context.Context, userID, addressID int64, shippingPackages string) (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":              ctx,
		"userID":           userID,
		"shippingPackages": shippingPackages,
	})

	_, err := u.addressUsecase.GetAddressByID(ctx, addressID)
	if err != nil {
		log.Error(err)
		return "", err
	}

	cart, err := u.cartUsecase.FindCartByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return "", err
	}

	user, err := u.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if user == nil {
		return "", ErrNotFound
	}

	price, err := u.cartUsecase.GetCartTotalPrice(ctx, userID)
	if err != nil {
		log.Error(err)
		return "", err
	}

	shipping, err := u.shippingUsecase.GetShippingPackageByServices(ctx, addressID, userID, shippingPackages)
	if err != nil {
		log.Error(err)
		return "", err
	}

	tx := u.db.Begin()

	userShipping := model.Shipping{
		AddressID:   addressID,
		Services:    shippingPackages,
		Description: shipping.Description,
		ETD:         shipping.ETD,
		Price:       shipping.Price,
	}
	shippingID, err := u.shippingUsecase.CreateShipping(ctx, tx, &userShipping)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return "", err
	}

	userShipping.ID = uint(shippingID)

	transID, err := u.transactionRepository.CreateTransaction(ctx, tx, &model.Transaction{
		UserID:     userID,
		Status:     model.TransactionStatusPending,
		Price:      price + shipping.Price,
		ShippingID: shippingID,
	})
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return "", err
	}

	err = u.transactionItemRepository.CreateMultipleItems(ctx, tx, transID, fromCartItemsToTransactionItems(transID, cart.Items))
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return "", err
	}

	err = u.cartItemRepository.DeleteByCartID(ctx, tx, userID)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return "", err
	}

	resp, _ := u.midtransSnapClient.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("%d", transID),
			GrossAmt: price + shipping.Price,
		},
		Items: fromCartItemsToMidtransItems(cart.Items, &userShipping),
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	})
	log.Info("transaction details = ", midtrans.TransactionDetails{
		OrderID:  fmt.Sprintf("%d", transID),
		GrossAmt: price + shipping.Price,
	})
	log.Info("transaction items = ", *fromCartItemsToMidtransItems(cart.Items, &userShipping))
	if len(resp.ErrorMessages) > 0 {
		log.Error(resp.ErrorMessages)
		tx.Rollback()
		return "", ErrCreateTransaction
	}

	err = u.transactionRepository.UpdateTransactionPaymentURL(ctx, tx, transID, resp.RedirectURL)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return "", err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error(err)
		return "", ErrCreateTransaction
	}

	return resp.RedirectURL, nil
}

func (u *transactionUsecaseImpl) UpdateTransactionStatus(ctx context.Context, transactionID int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":           ctx,
		"transactionID": transactionID,
	})

	trans, e := u.midtransCoreClient.CheckTransaction(fmt.Sprintf("%d", transactionID))
	if e != nil {
		log.Error(e)
		return e
	}

	status := model.ParseTransactionStatusFromString(trans.TransactionStatus)
	err := u.transactionRepository.UpdateTransactionStatus(ctx, transactionID, status)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *transactionUsecaseImpl) GetTransactionByUserID(ctx context.Context, userID int64) ([]*model.Transaction, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})

	transactions, err := u.transactionRepository.GetTransactionByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return transactions, nil
}

func (u *transactionUsecaseImpl) FindAllTransactions(ctx context.Context) ([]*model.Transaction, error) {
	log := logrus.WithField("ctx", ctx)

	trans, err := u.transactionRepository.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return trans, nil
}

func (u *transactionUsecaseImpl) FindByID(ctx context.Context, id int64, userID int64) (*model.Transaction, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"id":     id,
		"userID": userID,
	})

	user, err := u.userRepository.FindByID(ctx, userID)
	switch {
	case err != nil:
		log.Error(err)
		return nil, err
	case user == nil:
		return nil, ErrNotFound
	}

	trans, err := u.transactionRepository.FindByID(ctx, id)
	switch {
	case err != nil:
		log.Error(err)
		return nil, err
	case trans == nil:
		return nil, ErrNotFound
	case trans.UserID != userID && user.Role != model.RoleAdmin:
		return nil, ErrPermissionDenied
	}

	shipping, err := u.shippingUsecase.FindByID(ctx, trans.ShippingID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	trans.Shipping = shipping

	return trans, nil
}
