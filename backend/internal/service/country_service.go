package service

import (
	"context"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type CountryService interface {
	GetCountries(ctx context.Context) (*model.Countries, error)
}

type countryService struct {
	countryRepo repo.CountryRepo
}

func NewCountryService(countryRepo repo.CountryRepo) CountryService {
	return &countryService{
		countryRepo: countryRepo,
	}
}

func (c *countryService) GetCountries(ctx context.Context) (*model.Countries, error) {
	return c.countryRepo.Countries(ctx)
}
