package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
)

type addressUsecaseImpl struct {
	addressRepository model.AddressRepository
	validator         *validator.Validate
}

func NewAddressUsecase(addressRepository model.AddressRepository) model.AddressUsecase {
	return &addressUsecaseImpl{
		addressRepository: addressRepository,
		validator:         validator.New(),
	}
}

func (r *addressUsecaseImpl) GetProvinces(ctx context.Context) ([]*model.Province, error) {
	provinces, err := r.addressRepository.GetProvinces(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
		}).Error(err)
		return nil, err
	}

	return provinces, nil
}

func (r *addressUsecaseImpl) GetStatesByProvinceID(ctx context.Context, provinceID string) ([]*model.State, error) {
	states, err := r.addressRepository.GetStatesByProvinceID(ctx, provinceID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
		}).Error(err)
		return nil, err
	}

	return states, nil
}

func (u *addressUsecaseImpl) CreateAddress(ctx context.Context, address *model.Address) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":     ctx,
		"address": address,
	})

	err := u.validator.Struct(address)
	if err != nil {
		log.Error(err)
		return ErrValidation
	}

	addr, err := u.addressRepository.GetAddressByNameAndUserID(ctx, address.Name, address.UserID);
	if err != nil {
		log.Error(err)
		return err
	}
	if addr != nil {
		return ErrAddressNameAlreadyExist
	}

	state, err := u.addressRepository.GetStateByIDFromDB(ctx, address.StateID)
	if err != nil {
		log.Error(err)
		return err
	}
	if state == nil {
		state, err = u.addressRepository.GetStateByIDFromAPI(ctx, address.StateID)
		if err != nil {
			log.Error(err)
			return ErrNotFound
		}
		err = u.addressRepository.CreateState(ctx, state)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	province, err := u.addressRepository.GetProvinceByIDFromDB(ctx, address.ProvinceID)
	if err != nil {
		log.Error(err)
		return err
	}
	if province == nil {
		province, err = u.addressRepository.GetProvinceByIDFromAPI(ctx, address.ProvinceID)
		if err != nil {
			log.Error(err)
			return ErrNotFound
		}
		err = u.addressRepository.CreateProvince(ctx, province)
	}
	if err != nil {
		log.Error(err)
		return err
	}

	err = u.addressRepository.CreateAddress(ctx, address)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *addressUsecaseImpl) GetAddressesByUserID(ctx context.Context, userID int64) ([]*model.Address, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"userID": userID,
	})

	addresses, err := u.addressRepository.GetAddressesByUserID(ctx, userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return addresses, nil
}

func (u *addressUsecaseImpl) GetAddressByID(ctx context.Context, id int64) (*model.Address, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	address, err := u.addressRepository.GetAddressByID(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	logrus.Info("address in =", address)
	if address == nil {
		return nil, ErrNotFound
	}

	return address, nil
}
