package util

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoProvider struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoProvider(symmetricKey string) (TokenProvider, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoProvider{
		symmetricKey: []byte(symmetricKey),
		paseto:       paseto.NewV2(),
	}, nil
}

func (p *PasetoProvider) CreateToken(userID int64, role int, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, role, duration)
	if err != nil {
		return "", nil, err
	}
	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

func (p *PasetoProvider) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}
