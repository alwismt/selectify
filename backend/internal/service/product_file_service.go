package service

import (
	"context"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type productFileService struct {
	productFileRepo repo.ProductFileRepo
}

type ProductFileService interface {
	GetProductFileByProductID(ctx context.Context, proID uint) (*model.ProductFile, error)
}

func NewProductFileService(productFileRepo repo.ProductFileRepo) ProductFileService {
	return &productFileService{
		productFileRepo: productFileRepo,
	}
}

func (pfs *productFileService) GetProductFileByProductID(ctx context.Context, proID uint) (*model.ProductFile, error) {
	return pfs.productFileRepo.GetProductFileByProductID(ctx, proID)
}
