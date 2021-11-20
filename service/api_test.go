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
			Sign(ctx, email, 1).
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
			Sign(ctx, email, 1).
			Return("", jwt.ErrInvalidKey).Times(1)

		service := service.NewAPIService(userRepository, hasher, tokenSigner)

		accessToken, err := service.Login(ctx, email, password)
		assert.ErrorIs(t, err, jwt.ErrInvalidKey)
		assert.Empty(t, accessToken)
	})
}
