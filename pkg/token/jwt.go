package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type Claims struct {
	jwt.StandardClaims
	Email      string `json:"email"`
	MerchantID int    `json:"merchantID"`
}

type service struct {
	secret        string
	issuer        string
	signingMethod *jwt.SigningMethodHMAC
}

func NewJWTService(secret, issuer string) *service {
	return &service{secret: secret, issuer: issuer, signingMethod: jwt.SigningMethodHS256}
}

func (s *service) Sign(ctx context.Context, email string, merchantID int) (at string, err error) {
	const ops = "token.service.Sign"
	now := time.Now()

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	default:
		claims := Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: now.Add(12 * time.Hour).Unix(),
				IssuedAt:  now.Unix(),
				Issuer:    s.issuer,
				NotBefore: now.Unix(),
				Subject:   "access-token",
			},
			Email:      email,
			MerchantID: merchantID,
		}

		jwtToken := jwt.NewWithClaims(s.signingMethod, claims)
		signedToken, err := jwtToken.SignedString([]byte(s.secret))
		if err != nil {
			logger.Error(ctx, ops, "error trying to get JWT signed string: %v", err)
			return "", err
		}

		return signedToken, nil
	}
}

func (s *service) Extract(ctx context.Context, signedToken string) (jwt.MapClaims, error) {
	const ops = "token.service.Extract"

	select {
	case <-ctx.Done():
		return jwt.MapClaims{}, ctx.Err()
	default:
		token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
			if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("signing method invalid")
			} else if method != s.signingMethod {
				return nil, errors.New("signing method invalid")
			}
			return s.secret, nil
		})
		if err != nil {
			logger.Error(ctx, ops, "failed to parse signed token: %v", err)
			return nil, err
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.Error(ctx, ops, "token not valid")
			return jwt.MapClaims{}, errors.New("token not valid")
		}

		return claims, nil
	}
}
