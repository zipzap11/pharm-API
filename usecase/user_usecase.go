package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/config"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
	"gorm.io/gorm"
)

type userUsecaseImpl struct {
	userRepository    model.UserRepository
	validator         *validator.Validate
	tokenProvider     util.TokenProvider
	sessionRepository model.SessionRepository
	db                *gorm.DB
	cartRepository    model.CartRepository
}

func NewUserUsecase(userRepository model.UserRepository, validator *validator.Validate, tokenProvider util.TokenProvider, sessionRepository model.SessionRepository, db *gorm.DB, cartRepository model.CartRepository) model.UserUsecase {
	return &userUsecaseImpl{
		userRepository:    userRepository,
		validator:         validator,
		tokenProvider:     tokenProvider,
		sessionRepository: sessionRepository,
		db:                db,
		cartRepository:    cartRepository,
	}
}

func (u *userUsecaseImpl) CreateUser(ctx context.Context, user *model.User) error {
	log := logrus.WithField("user", user)
	err := u.validator.Struct(user)
	if err != nil {
		log.Error(err)
		return err
	}

	usr, err := u.userRepository.FindUserByEmail(ctx, user.Email)
	switch {
	case err != nil:
		log.Error(err)
		return err
	case usr != nil:
		return ErrEmailAlreadyUsed
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error(err)
		return err
	}

	user.Password = hashedPassword
	user.Role = model.RoleUser

	tx := u.db.Begin()
	userID, err := u.userRepository.CreateUser(ctx, tx, user)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}

	err = u.cartRepository.CreateCart(ctx, tx, userID)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *userUsecaseImpl) Login(ctx context.Context, email string, password string) (string, string, error) {
	log := logrus.WithFields(logrus.Fields{
		"email":    email,
		"password": password,
		"ctx":      ctx,
	})

	err := u.validator.Var(email, "email,required")
	if err != nil {
		log.Error(err)
		return "", "", ErrInvalidEmail
	}
	err = u.validator.Var(password, "min=6")
	if err != nil {
		log.Error(err)
		return "", "", ErrInvalidEmail
	}

	user, err := u.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	if user == nil {
		return "", "", ErrNotFound
	}

	err = util.CheckPassword(password, user.Password)
	if err != nil {
		return "", "", ErrInvalidCredential
	}

	log.Info("user role =", user.Role, " id =", user.ID)
	accessToken, _, err := u.tokenProvider.CreateToken(int64(user.ID), int(user.Role), config.GetAccessTokenDuration())
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	refreshToken, refreshPayload, err := u.tokenProvider.CreateToken(int64(user.ID), int(user.Role), config.GetRefreshTokenDuration())
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	err = u.sessionRepository.CreateSession(ctx, &model.Session{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		UserID:       refreshPayload.UserID,
		Role:         user.Role,
		CreatedAt:    refreshPayload.CreatedAt,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		log.Error(err)
	}

	return accessToken, refreshToken, nil
}

func (u *userUsecaseImpl) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := u.userRepository.FindUserByID(ctx, id)
	switch {
	case err != nil:
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
			"id":  id,
		}).Error(err)
		return nil, err
	case user == nil:
		return nil, ErrNotFound
	default:
		return user, nil
	}
}

func (u *userUsecaseImpl) CreateSuperUser(ctx context.Context, user *model.User) error {
	log := logrus.WithField("user", user)
	err := u.validator.Struct(user)
	if err != nil {
		log.Error(err)
		return err
	}

	usr, err := u.userRepository.FindUserByEmail(ctx, user.Email)
	switch {
	case err != nil:
		log.Error(err)
		return err
	case usr != nil:
		return ErrEmailAlreadyUsed
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error(err)
		return err
	}

	user.Password = hashedPassword
	user.Role = model.RoleAdmin

	tx := u.db.Begin()
	_, err = u.userRepository.CreateUser(ctx, tx, user)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *userUsecaseImpl) FindAllUsers(ctx context.Context) ([]*model.User,error) {
	log := logrus.WithField("ctx", ctx)
	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return users, err
}