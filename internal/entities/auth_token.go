package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthTokenPayload struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewAuthTokenPayload(userId string, duration time.Duration) (*AuthTokenPayload, error) {
	if userId == "" {
		return nil, errors.New("empty username")
	}

	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &AuthTokenPayload{
		ID:        tokenID.String(),
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func CreateAuthTokenPayload(id, userId string, issuedAt, expiredAt time.Time) (*AuthTokenPayload, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	if userId == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if issuedAt.IsZero() {
		return nil, fmt.Errorf("issued at is required")
	}

	if expiredAt.IsZero() {
		return nil, fmt.Errorf("expired at is required")
	}

	return &AuthTokenPayload{
		ID:        uuid.New().String(),
		UserId:    userId,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}

func (p *AuthTokenPayload) IsExpired() bool {
	return time.Now().After(p.ExpiredAt)
}
