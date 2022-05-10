package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zipzap11/pharm-API/model"
	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) model.SessionRepository {
	return &SessionRepositoryImpl{
		db: db,
	}
}

func (r *SessionRepositoryImpl) CreateSession(ctx context.Context, session *model.Session) error {
	err := r.db.Create(session).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":     ctx,
			"session": session,
		})
		return err
	}
	return nil
}

func (r *SessionRepositoryImpl) GetSessionByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	var session model.Session
	err := r.db.Model(&model.Session{}).Where("id = ?", id.String()).First(&session).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		logrus.WithFields(logrus.Fields{
			"ctx": ctx,
			"id":  id,
		}).Error(err)
		return nil, err
	}

	return &session, nil
}
