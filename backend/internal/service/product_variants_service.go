package service

import (
	"context"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type ProductVariantsService interface {
	GetProVariantsForProduct(ctx context.Context, product *model.Product) (model.ProductVariants, error)
}

type productVariantsService struct {
	proVariantsRepo repo.ProductVariantsRepo
	//productFileRepo repo.ProductFileRepo
}

func NewProductVariantsService(proVariantsRepo repo.ProductVariantsRepo, productFileRepo repo.ProductFileRepo) ProductVariantsService {
	return &productVariantsService{
		proVariantsRepo: proVariantsRepo,
		//productFileRepo: productFileRepo,
	}
}

func (pvs *productVariantsService) GetProVariantsForProduct(ctx context.Context, product *model.Product) (model.ProductVariants, error) {
	return pvs.proVariantsRepo.GetVariantsForProduct(ctx, product.ID)
}
