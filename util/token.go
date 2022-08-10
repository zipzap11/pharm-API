package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *Payload) Valid() error {
	fmt.Println("expired at = ", p.ExpiredAt)
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(userID int64, role int, duration time.Duration) (*Payload, error) {
	fmt.Println("time now = ", time.Now())
	fmt.Println("add = ", time.Now().Add(duration))
	fmt.Println("duration = ", duration)
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        id,
		UserID:    userID,
		Role:      role,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

type TokenProvider interface {
	CreateToken(UserID int64, role int, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
