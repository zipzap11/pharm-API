package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const AuthorizationPayloadKey = "authorization_payload"

type Role int

var (
	RoleUser  Role = 1
	RoleAdmin Role = 0
)

type Session struct {
	ID           uuid.UUID
	RefreshToken string
	UserID       int64
	Role         Role
	IsBlocked    bool
	CreatedAt    time.Time
	ExpiredAt    time.Time
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
}

type SessionUsecase interface {
	RefreshSession(ctx context.Context, token string) (string, string, error)
	CheckSession(ctx context.Context, token string) error
}
