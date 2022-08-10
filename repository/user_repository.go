package repository

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, tx *gorm.DB, user *model.User) (int64, error) {
	err := tx.Create(user).Error
	if err != nil {
		logrus.WithField("user", user).Error(err)
		return 0, err
	}
	return int64(user.ID), nil
}

func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	var user *model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context) ([]*model.User, error) {
	log := logrus.WithField("ctx", ctx)

	var users []*model.User
	if err := r.db.Find(&users).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id": id,
	})
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}

	return &user, nil
}