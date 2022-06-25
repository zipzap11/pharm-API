package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type cartUsecaseImpl struct {
	cartRepository     model.CartRepository
	cartItemRepository model.CartItemRepository
}

func NewCartUsecase(cartRepository model.CartRepository, cartItemRepository model.CartItemRepository) model.CartUsecase {
	return &cartUsecaseImpl{
		cartRepository:     cartRepository,
		cartItemRepository: cartItemRepository,
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
	cart.TotalPrice = u.getTotalPriceByItems(cart.Items)
	return cart, nil
}

func (u *cartUsecaseImpl) AddItemToCart(ctx context.Context, userID, productID int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"userID":    userID,
		"productID": productID,
	})
	cart, err := u.cartRepository.FindCartByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return err
	}

	itemID, err := u.cartItemRepository.FindItemIDByCartIDAndProductID(ctx, int64(cart.ID), productID)
	if err != nil {
		log.Error(err)
		return err
	}
	if itemID != 0 {
		return ErrItemAlreadyExist
	}

	err = u.cartItemRepository.CreateItem(ctx, int64(cart.ID), productID)
	if err != nil {
		logrus.WithField("cartID", cart.ID).Error(err)
	}

	return err
}

func (u *cartUsecaseImpl) UpdateItemQuantity(ctx context.Context, itemID, userID, quantity int64, opsType model.QuantityUpdateType) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":      ctx,
		"itemID":   itemID,
		"userID":   userID,
		"quantity": quantity,
	})

	cartID, err := u.cartRepository.GetCartIDByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return err
	}
	if cartID == 0 {
		return ErrNotFound
	}
	item, err := u.cartItemRepository.FindItemByID(ctx, itemID)
	log.Info("cart id =", cartID)
	log.Info("item.CartID =", item.CartID)
	switch {
	case err != nil:
		log.Error(err)
		return err
	case item == nil:
		return ErrNotFound
	case item.CartID != cartID:
		return ErrPermissionDenied
	// check if item's quantity will be zero after update
	case quantity == 0 && opsType == model.UPDATE_QUANTITY_OPERATION_TYPE,
		item.Quantity == 1 && opsType == model.SUBTRACT_QUANTITY_OPERATION_TYPE:
		err = u.cartItemRepository.DeleteItemByID(ctx, int64(item.ID))
	case opsType == model.UPDATE_QUANTITY_OPERATION_TYPE:
		err = u.cartItemRepository.UpdateItemQuantity(ctx, int64(item.ID), quantity)
	case opsType == model.ADD_QUANTITY_OPERATION_TYPE,
		opsType == model.SUBTRACT_QUANTITY_OPERATION_TYPE:
		err = u.cartItemRepository.UpdateItemQuantityByType(ctx, int64(item.ID), opsType)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *cartUsecaseImpl) RemoveItemFromCart(ctx context.Context, itemID, userID int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
		"itemID": itemID,
	})
	cartID, err := u.cartRepository.GetCartIDByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return err
	}
	if cartID == 0 {
		return ErrNotFound
	}

	item, err := u.cartItemRepository.FindItemByID(ctx, itemID)
	if err != nil {
		log.Error(err)
		return err
	}
	if item == nil {
		return ErrNotFound
	}
	if item.CartID != cartID {
		return ErrPermissionDenied
	}

	err = u.cartItemRepository.DeleteItemByID(ctx, itemID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *cartUsecaseImpl) GetCartTotalPrice(ctx context.Context, userID int64) (int64, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})
	cartID, err := u.cartRepository.GetCartIDByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	if cartID == 0 {
		return 0, ErrNotFound
	}

	items, err := u.cartItemRepository.GetItemsByCartID(ctx, cartID)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	var total int64
	for _, v := range items {
		total += (v.Product.Price * v.Quantity)
	}

	return total, nil
}

func (u *cartUsecaseImpl) getTotalPriceByItems(items []*model.CartItem) int64 {
	var total int64
	for _, v := range items {
		total += (v.Product.Price * v.Quantity)
	}

	return total
}