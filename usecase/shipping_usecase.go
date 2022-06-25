package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type shippingUsecaseImpl struct {
	shippingRepository model.ShippingRepository
	cartItemRepository model.CartItemRepository
	cartRepository     model.CartRepository
	addressRepository  model.AddressRepository
}

func NewShippingUsecase(
	shippingRepository model.ShippingRepository,
	cartItemRepository model.CartItemRepository,
	cartRepository model.CartRepository,
	addressRepository model.AddressRepository,
) model.ShippingUsecase {
	return &shippingUsecaseImpl{
		shippingRepository: shippingRepository,
		cartItemRepository: cartItemRepository,
		cartRepository:     cartRepository,
		addressRepository:  addressRepository,
	}
}

func (u *shippingUsecaseImpl) GetShippingPackages(ctx context.Context, addressID, userID int64) ([]*model.Shipping, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":     ctx,
		"stateID": addressID,
	})

	items, err := u.cartItemRepository.GetItemsByCartID(ctx, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range items {
		log.Infof("item - %d weight = %v\n", v.ID,v.Product.Weight)
	}

	address, err := u.addressRepository.GetAddressByID(ctx, addressID)
	switch {
	case err != nil:
		log.Error(err)
		return nil, err
	case address == nil:
		return nil, ErrNotFound
	case address.UserID != userID:
		return nil, ErrInvalidAddress
	}

	weight := u.countItemsWeight(items)
	if weight < 1 {
		weight = 1
	}
	shippings, err := u.shippingRepository.GetShippingsPackages(ctx, address.StateID, weight)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return shippings, nil
}

func (u *shippingUsecaseImpl) GetShippingPackageByServices(ctx context.Context, addressID, userID int64, service string) (*model.Shipping, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"addressID": addressID,
		"userID":    userID,
		"service":   service,
	})

	pkgs, err := u.GetShippingPackages(ctx, addressID, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, v := range pkgs {
		if v.Services == service {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

func (u *shippingUsecaseImpl) CreateShipping(ctx context.Context, tx *gorm.DB, shipping *model.Shipping) (int64, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":      ctx,
		"shipping": shipping,
	})

	id, err := u.shippingRepository.CreateShipping(ctx, tx, shipping)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return id, nil
}

func (u *shippingUsecaseImpl) countItemsWeight(items []*model.CartItem) float64 {
	var weight float64
	for _, v := range items {
		weight += (v.Product.Weight * float64(v.Quantity))
	}
	return weight
}
