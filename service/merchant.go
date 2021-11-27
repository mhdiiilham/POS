package service

import (
	"context"
	"time"

	"github.com/mhdiiilham/POS/entity/merchant"
	"github.com/mhdiiilham/POS/pkg/logger"
)

type merchantService struct {
	merchantRepository merchant.Repository
}

func NewMerchantService(merchantRepository merchant.Repository) *merchantService {
	return &merchantService{merchantRepository: merchantRepository}
}

func (s *merchantService) NewMerchant(ctx context.Context, entity merchant.Merchant) (merchant.Merchant, error) {
	const ops = "service.merchantService.NewMerchant"
	now := time.Now()

	id, err := s.merchantRepository.Create(ctx, entity)
	if err != nil {
		logger.Error(ctx, ops, "unexpected error: %v", err)
		return entity, err
	}

	entity.ID = id
	entity.CreatedAt = now
	entity.UpdatedAt = now
	return entity, nil
}
