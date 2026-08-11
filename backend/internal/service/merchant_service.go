package service

import (
	"context"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type MerchantService interface {
	GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error)
}

type merchantService struct {
	merchantRepo repo.MerchantRepo
}

func NewMerchantService(merchantRepo repo.MerchantRepo) MerchantService {
	return &merchantService{
		merchantRepo: merchantRepo,
	}

}

func (s *merchantService) GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error) {
	return s.merchantRepo.GetMerchant(ctx, merchantID)
}
