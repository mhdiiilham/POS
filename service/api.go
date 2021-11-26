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

	accessToken, err = s.tokenSigner.Sign(ctx, entity.ID, entity.Email, entity.MerchantID)
	if err != nil {
		logger.Error(ctx, ops, "error trying to compare password %v", err)
		return "", err
	}

	return accessToken, nil
}

func (s *apiService) CreateUser(ctx context.Context, entity user.User) (userID int, err error) {
	const ops = "service.apiService.CreateUser"
	var hashedPwd string
	var insertedID int64
	var u *user.User

	if entity.Email == "" || entity.FirstName == "" || len(entity.Password) < 8 {
		return 0, user.ErrInvalidCreateParameters
	}

	hashedPwd, err = s.hasher.HashPassword(ctx, entity.Password)
	if err != nil {
		logger.Error(ctx, ops, "error when trying to hash password: %v", err)
		return 0, err
	}

	u, err = s.userRepository.FindUserByEmail(ctx, entity.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error(ctx, ops, "unexpected error happened %v", err)
		return 0, err
	}

	if u != nil {
		return 0, user.ErrEmailNotUnique
	}

	entity.Password = hashedPwd
	insertedID, err = s.userRepository.Create(ctx, entity)
	if err != nil {
		logger.Error(ctx, ops, "error when trying to insert entity to db: %v", err)
		return 0, err
	}
	return int(insertedID), nil
}

func (s *apiService) GetUsers(ctx context.Context, merchantID, lastID, limit int) (users []user.User, totalData int, err error) {
	const ops = "service.apiService.GetUsers"
	paginationOpts := user.RepositoryGetUserPaginationOptions{
		Limit:  limit,
		Cursor: lastID,
	}

	users, totalData, err = s.userRepository.Get(ctx, merchantID, &paginationOpts)
	if err != nil {
		logger.Error(ctx, ops, "unexpected error %v", err)
		return
	}

	return
}

func (s *apiService) DeleteUser(ctx context.Context, userID int) error {
	const ops = "service.apiService.DeleteUser"
	err := s.userRepository.Remove(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		}

		logger.Error(ctx, ops, "error removing user %v", err)
		return err
	}

	return nil
}
