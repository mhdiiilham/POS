package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type apiService struct {
	userRepository user.Repository
	hasher         Hasher
	tokenSigner    TokenSigner
}

func NewAPIService(userRepository user.Repository, pwdHasher Hasher, tokenSigner TokenSigner) *apiService {
	return &apiService{
		userRepository: userRepository,
		hasher:         pwdHasher,
		tokenSigner:    tokenSigner,
	}
}

func (s *apiService) Login(ctx context.Context, email, password string) (accessToken string, err error) {
	const ops = "service.user.Login"
	entity, err := s.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", user.ErrInvalidEmailAndPasword
		}

		logger.Error(ctx, ops, "error trying to find user by email: %v", err)
		return "", err
	}

	if err := s.hasher.ComparePassword(ctx, entity.Password, password); err != nil {
		logger.Error(ctx, ops, "error trying to compare password %v", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", user.ErrInvalidEmailAndPasword
		}
		return "", err
	}

	accessToken, err = s.tokenSigner.Sign(ctx, entity.Email, entity.MerchantID)
	if err != nil {
		logger.Error(ctx, ops, "error trying to compare password %v", err)
		return "", err
	}

	return accessToken, nil
}
