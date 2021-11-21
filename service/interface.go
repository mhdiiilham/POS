package service

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type Hasher interface {
	HashPassword(ctx context.Context, password string) (hashed string, err error)
	ComparePassword(ctx context.Context, hashedPassword, password string) error
}

type TokenSigner interface {
	Sign(ctx context.Context, userID int, email string, merchantID int) (at string, err error)
	Extract(ctx context.Context, signedToken string) (jwt.MapClaims, error)
}
