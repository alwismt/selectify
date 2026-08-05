package service

import (
	"context"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type productService struct {
	productRepo     repo.ProductRepo
	productFileRepo repo.ProductFileRepo
}

type ProductService interface {
	GetProductById(ctx context.Context, id uint) (*model.Product, error)
	GetProductFileByProductID(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*model.Product, error)
	GetProducts(ctx context.Context) (model.Products, error)
}

func NewProductService(productRepo repo.ProductRepo, fileRepo repo.ProductFileRepo) ProductService {
	return &productService{
		productRepo:     productRepo,
		productFileRepo: fileRepo,
	}
}

func (ps *productService) GetProductById(ctx context.Context, id uint) (*model.Product, error) {
	product, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, logger.Error(ctx, err, "error getting product by id")
	}

	return ps.GetProductFileByProductID(ctx, product)
}

func (ps *productService) GetProductFileByProductID(ctx context.Context, product *model.Product) (*model.Product, error) {
	var err error = nil
	product.ProductFile, err = ps.productFileRepo.GetProductFileByProductID(ctx, product.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "error getting product file by product id")
	}
	return product, nil
}

func (ps *productService) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	return ps.productRepo.GetProductBySlug(ctx, slug)
}

func (ps *productService) GetProducts(ctx context.Context) (model.Products, error) {
	return ps.productRepo.GetProducts(ctx)
}
