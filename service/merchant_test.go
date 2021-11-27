package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/mhdiiilham/POS/entity/merchant"
	"github.com/mhdiiilham/POS/entity/merchant/mock"
	"github.com/mhdiiilham/POS/service"
	"github.com/stretchr/testify/assert"
)

func Test_merchantService_NewMerchant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("failed - inserting to db", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		payload := merchant.Merchant{Name: faker.Name()}
		merchantRepository := mock.NewMockRepository(ctrl)

		merchantRepository.
			EXPECT().
			Create(ctx, payload).
			Return(0, sql.ErrConnDone).
			Times(1)

		s := service.NewMerchantService(merchantRepository)
		_, err := s.NewMerchant(ctx, payload)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		payload := merchant.Merchant{Name: faker.Name()}
		merchantRepository := mock.NewMockRepository(ctrl)

		merchantRepository.
			EXPECT().
			Create(ctx, payload).
			Return(1, nil).
			Times(1)

		s := service.NewMerchantService(merchantRepository)
		resp, err := s.NewMerchant(ctx, payload)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.CreatedAt)
		assert.NotEmpty(t, resp.UpdatedAt)
		assert.Equal(t, resp.ID, 1)
	})
}
