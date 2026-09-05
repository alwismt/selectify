package service

import (
	"context"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type SiteConfigService interface {
	GetSiteConfig(ctx context.Context) (*model.SiteConfig, error)
}

type siteConfigService struct {
	currencyRepo repo.CurrencyRepo
}

func NewSiteConfigService(currencyRepo repo.CurrencyRepo) SiteConfigService {
	return &siteConfigService{
		currencyRepo: currencyRepo,
	}
}

func (s *siteConfigService) GetSiteConfig(ctx context.Context) (*model.SiteConfig, error) {
	config := new(model.SiteConfig)

	currency, err := s.currencyRepo.GetDefaultCurrency(ctx)
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to get default currency")
	}

	config.Currency = currency

	return config, nil
}
