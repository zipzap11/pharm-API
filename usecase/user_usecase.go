package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/config"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type userUsecaseImpl struct {
	userRepository    model.UserRepository
	validator         *validator.Validate
	tokenProvider     util.TokenProvider
	sessionRepository model.SessionRepository
}

func NewUserUsecase(userRepository model.UserRepository, validator *validator.Validate, tokenProvider util.TokenProvider, sessionRepository model.SessionRepository) model.UserUsecase {
	return &userUsecaseImpl{
		userRepository:    userRepository,
		validator:         validator,
		tokenProvider:     tokenProvider,
		sessionRepository: sessionRepository,
	}
}

func (u *userUsecaseImpl) CreateUser(ctx context.Context, user *model.User) error {
	log := logrus.WithField("user", user)
	err := u.validator.Struct(user)
	if err != nil {
		log.Error(err)
		return err
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Error(err)
		return err
	}

	user.Password = hashedPassword
	err = u.userRepository.CreateUser(ctx, user)
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

	accessToken, _, err := u.tokenProvider.CreateToken(int64(user.ID), config.GetAccessTokenDuration())
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	refreshToken, refreshPayload, err := u.tokenProvider.CreateToken(int64(user.ID), config.GetRefreshTokenDuration())
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	err = u.sessionRepository.CreateSession(ctx, &model.Session{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		UserID:       refreshPayload.UserID,
		CreatedAt:    refreshPayload.CreatedAt,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		log.Error(err)
	}

	return accessToken, refreshToken, nil
}
