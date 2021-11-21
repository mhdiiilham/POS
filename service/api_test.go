package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/mhdiiilham/POS/entity/user"
	"github.com/mhdiiilham/POS/entity/user/mock"
	"github.com/mhdiiilham/POS/service"
	smock "github.com/mhdiiilham/POS/service/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_apiService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		hashedPassword := faker.Password()
		jwt := faker.Jwt()

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(&user.User{
				ID:         1,
				Email:      email,
				Password:   hashedPassword,
				MerchantID: 1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  nil,
			}, nil).Times(1)

		hasher.EXPECT().
			ComparePassword(ctx, hashedPassword, password).
			Return(nil).
			Times(1)

		tokenSigner.
			EXPECT().
			Sign(ctx, 1, email, 1).
			Return(jwt, nil).Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.Equal(t, jwt, accessToken)
	})

	t.Run("fail - user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		expectedErr := user.ErrInvalidEmailAndPasword

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(nil, sql.ErrNoRows).
			Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, accessToken)
	})

	t.Run("fail - unkown error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		expectedErr := sql.ErrConnDone

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(nil, sql.ErrConnDone).
			Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, expectedErr)
		assert.Empty(t, accessToken)
	})

	t.Run("failed - wrong password", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		hashedPassword := faker.Password()

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(&user.User{
				ID:         1,
				Email:      email,
				Password:   hashedPassword,
				MerchantID: 1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  nil,
			}, nil).Times(1)

		hasher.EXPECT().
			ComparePassword(ctx, hashedPassword, password).
			Return(bcrypt.ErrMismatchedHashAndPassword).
			Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, user.ErrInvalidEmailAndPasword)
		assert.Empty(t, accessToken)
	})

	t.Run("failed - bcyrpt unknown error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		hashedPassword := faker.Password()

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(&user.User{
				ID:         1,
				Email:      email,
				Password:   hashedPassword,
				MerchantID: 1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  nil,
			}, nil).Times(1)

		hasher.EXPECT().
			ComparePassword(ctx, hashedPassword, password).
			Return(bcrypt.ErrHashTooShort).
			Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, bcrypt.ErrHashTooShort)
		assert.Empty(t, accessToken)
	})

	t.Run("failed - sign token error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		email := faker.Email()
		password := faker.Password()
		hashedPassword := faker.Password()

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, email).
			Return(&user.User{
				ID:         1,
				Email:      email,
				Password:   hashedPassword,
				MerchantID: 1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  nil,
			}, nil).Times(1)

		hasher.EXPECT().
			ComparePassword(ctx, hashedPassword, password).
			Return(nil).
			Times(1)

		tokenSigner.
			EXPECT().
			Sign(ctx, 1, email, 1).
			Return("", jwt.ErrInvalidKey).Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, jwt.ErrInvalidKey)
		assert.Empty(t, accessToken)
	})
}

func Test_apiService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("failed - required payload is empty", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		payload := user.User{}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.Empty(t, resp)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, user.ErrInvalidCreateParameters)
	})

	t.Run("failed - hashing password", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		lastname := faker.LastName()
		password := faker.Password()
		payload := user.User{
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  &lastname,
			Password:  password,
		}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		hasher.
			EXPECT().
			HashPassword(ctx, password).
			Return("", bcrypt.ErrHashTooShort).
			Times(1)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.Empty(t, resp)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, bcrypt.ErrHashTooShort)
	})

	t.Run("failed - inserting to DB", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		lastname := faker.LastName()
		password := faker.Password()
		payload := user.User{
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  &lastname,
			Password:  password,
		}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		hasher.
			EXPECT().
			HashPassword(ctx, password).
			Return(password, nil).
			Times(1)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, payload.Email).
			Return(nil, sql.ErrNoRows).
			Times(1)

		userRepository.
			EXPECT().
			Create(ctx, payload).
			Return(int64(0), sql.ErrConnDone).
			Times(1)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.Empty(t, resp)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, sql.ErrConnDone)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		lastname := faker.LastName()
		password := faker.Password()
		payload := user.User{
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  &lastname,
			Password:  password,
		}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		hasher.
			EXPECT().
			HashPassword(ctx, password).
			Return(password, nil).
			Times(1)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, payload.Email).
			Return(nil, sql.ErrNoRows).
			Times(1)

		userRepository.
			EXPECT().
			Create(ctx, payload).
			Return(int64(1), nil).
			Times(1)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp)
	})

	t.Run("email not unique", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		lastname := faker.LastName()
		password := faker.Password()
		payload := user.User{
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  &lastname,
			Password:  password,
		}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		hasher.
			EXPECT().
			HashPassword(ctx, password).
			Return(password, nil).
			Times(1)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, payload.Email).
			Return(&user.User{}, sql.ErrNoRows).
			Times(1)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.Empty(t, resp)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, user.ErrEmailNotUnique)
	})

	t.Run("failed - unkownn when checking if user unique", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		lastname := faker.LastName()
		password := faker.Password()
		payload := user.User{
			Email:     faker.Email(),
			FirstName: faker.FirstName(),
			LastName:  &lastname,
			Password:  password,
		}

		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)

		hasher.
			EXPECT().
			HashPassword(ctx, password).
			Return(password, nil).
			Times(1)

		userRepository.
			EXPECT().
			FindUserByEmail(ctx, payload.Email).
			Return(&user.User{}, sql.ErrConnDone).
			Times(1)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		resp, err := s.CreateUser(ctx, payload)
		assert.Empty(t, resp)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, sql.ErrConnDone)
	})
}

func Test_apiService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("failed - db error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)
		merchantID := 1
		expectedErr := sql.ErrConnDone

		opts := user.RepositoryGetUserPaginationOptions{
			Limit:  10,
			Cursor: 0,
		}

		userRepository.
			EXPECT().
			Get(ctx, merchantID, &opts).
			Return([]user.User{}, 0, sql.ErrConnDone)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		users, totalData, err := s.GetUsers(ctx, merchantID, opts.Cursor, opts.Limit)
		assert.Empty(t, users)
		assert.Empty(t, totalData)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userRepository := mock.NewMockRepository(ctrl)
		hasher := smock.NewMockHasher(ctrl)
		tokenSigner := smock.NewMockTokenSigner(ctrl)
		merchantID := 1

		opts := user.RepositoryGetUserPaginationOptions{
			Limit:  10,
			Cursor: 0,
		}

		userRepository.
			EXPECT().
			Get(ctx, merchantID, &opts).
			Return([]user.User{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}, 1764, nil)

		s := service.NewAPIService(userRepository, hasher, tokenSigner)
		users, totalData, err := s.GetUsers(ctx, merchantID, opts.Cursor, opts.Limit)
		assert.NoError(t, err)
		assert.Equal(t, totalData, 1764)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 10)
	})
}
