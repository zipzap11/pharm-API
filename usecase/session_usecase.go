package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/config"
	"github.com/zipzap11/pharm-API/model"
	"github.com/zipzap11/pharm-API/util"
)

type SessionUsecaseImpl struct {
	sessionRepository model.SessionRepository
	tokenProvider     util.TokenProvider
}

func NewSessionUsecase(sessionRepository model.SessionRepository, tokenProvider util.TokenProvider) model.SessionUsecase {
	return &SessionUsecaseImpl{
		sessionRepository: sessionRepository,
		tokenProvider:     tokenProvider,
	}
}

func (u *SessionUsecaseImpl) RefreshSession(ctx context.Context, token string) (string, string, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":   ctx,
		"token": token,
	})
	payload, err := u.tokenProvider.VerifyToken(token)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	session, err := u.sessionRepository.GetSessionByID(ctx, payload.ID)
	if err != nil {
		log.WithField("payload", payload)
		return "", "", err
	}
	if session == nil {
		return "", "", ErrNotFound
	}
	if session.IsBlocked {
		return "", "", ErrBlockedSession
	}
	if payload.UserID != session.UserID {
		return "", "", ErrIncorrectUserToken
	}
	if session.RefreshToken != token {
		return "", "", ErrMissmatchedToken
	}

	accessToken, _, err := u.tokenProvider.CreateToken(session.UserID, config.GetAccessTokenDuration())
	if err != nil {
		log.WithField("session", session).Error(err)
		return "", "", err
	}

	refreshToken, refreshPayload, err := u.tokenProvider.CreateToken(session.UserID, config.GetRefreshTokenDuration())
	if err != nil {
		log.WithField("session", session).Error(err)
		return "", "", err
	}

	err = u.sessionRepository.CreateSession(ctx, &model.Session{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		UserID:       payload.UserID,
		CreatedAt:    payload.CreatedAt,
		ExpiredAt:    payload.ExpiredAt,
	})
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *SessionUsecaseImpl) CheckSession(ctx context.Context, token string) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":   ctx,
		"token": token,
	})
	_, err := u.tokenProvider.VerifyToken(token)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
