package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type MerchantService interface {
	GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant *model.Merchant, req *request.UpdateMerchantRequest) error
	UpdateMerchantImage(ctx context.Context, merchant *model.Merchant, fileStream model.FileStream) error
}

type merchantService struct {
	storageService StorageService

	countryRepo  repo.CountryRepo
	merchantRepo repo.MerchantRepo
	txRepo       repo.TransactionRepo
}

func NewMerchantService(storageService StorageService, countryRepo repo.CountryRepo, merchantRepo repo.MerchantRepo,
	txRepo repo.TransactionRepo) MerchantService {
	return &merchantService{
		storageService: storageService,
		countryRepo:    countryRepo,
		merchantRepo:   merchantRepo,
		txRepo:         txRepo,
	}

}

func (s *merchantService) GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error) {
	return s.merchantRepo.GetMerchant(ctx, merchantID)
}

func (s *merchantService) UpdateMerchant(ctx context.Context, merchant *model.Merchant, req *request.UpdateMerchantRequest) error {
	if req == nil {
		return logger.Error(ctx, nil, "Update merchant request cannot be nil")
	}
	if req.Name != nil {
		merchant.Name = *req.Name
	}
	if req.Description != nil {
		merchant.Description = req.Description
	}
	if req.CountryCode != nil {
		if !s.countryRepo.CountryExists(ctx, *req.CountryCode) {
			_ = logger.Error(ctx, nil, httpx.ErrInvalidCountry.Error())
			return httpx.ErrInvalidCountry
		}
		merchant.CountryCode = *req.CountryCode
	}

	return s.merchantRepo.UpdateMerchant(ctx, merchant)
}

func (s *merchantService) UpdateMerchantImage(ctx context.Context, merchant *model.Merchant, fileStream model.FileStream) error {
	oldLogo := merchant.Logo
	tx, err := s.txRepo.BeginTransaction(ctx)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to begin transaction")
	}
	defer tx.End()

	newLogo := imageObjectKey(uuid.New().String(), merchant.Slug)
	merchant.Logo = &newLogo
	if err = s.merchantRepo.UpdateMerchantWithTx(ctx, merchant, tx.Transaction); err != nil {
		return logger.Errorf(ctx, err, "failed to update merchant")
	}

	if err = s.storageService.UploadFile(ctx, fileStream, newLogo); err != nil {
		return logger.Errorf(ctx, err, "failed to upload merchant logo")
	}

	tx.CanCommit = true

	if oldLogo != nil {
		go func(oldLogo string) {
			if err = s.storageService.DeleteFile(context.Background(), oldLogo); err != nil {
				_ = logger.Errorf(ctx, err, "failed to delete merchant logo for object %s", oldLogo)
			}
		}(*oldLogo)
	}
	return nil
}

func imageObjectKey(key, name string) string {
	return fmt.Sprintf("merchant/logo/%s-%s", key, name)
}
