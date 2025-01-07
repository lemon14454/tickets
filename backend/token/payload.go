package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("token is invalid")
var ErrExpiredToken = errors.New("token has expired")

type Payload struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID.String(),
		},
	}

	return payload, nil
}

func isPayloadExpired(payload jwt.Claims) error {
	t, err := payload.GetExpirationTime()
	if err != nil {
		return err
	}
	if time.Now().After(t.Time) {
		return ErrExpiredToken
	}
	return nil
}
