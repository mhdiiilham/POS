package hasher

import (
	"context"

	"github.com/mhdiiilham/POS/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(ctx context.Context, password string) (hashed string, err error)
	ComparePassword(ctx context.Context, hashedPassword, password string) error
}

type hasher struct {
}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) HashPassword(ctx context.Context, password string) (string, error) {
	const ops = "pkg.hasher.HashPassword"

	select {
	case <-ctx.Done():
		logger.Info(ctx, ops, ctx.Err().Error())
		return "", ctx.Err()

	default:
		p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			logger.Error(ctx, ops, "error trying to hash password: %v", err)
			return "", err
		}
		return string(p), nil
	}
}

func (h *hasher) ComparePassword(ctx context.Context, hashedPassword, password string) error {
	const ops = "pkg.hasher.ComparePassword"

	select {
	case <-ctx.Done():
		logger.Info(ctx, ops, ctx.Err().Error())
		return ctx.Err()

	default:
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	}
}
