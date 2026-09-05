package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"

	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type productService struct {
	storageService  StorageService
	categoryService CategoryService

	productRepo     repo.ProductRepo
	productFileRepo repo.ProductFileRepo
	txRepo          repo.TransactionRepo
}

type ProductService interface {
	GetProductById(ctx context.Context, id uint) (*model.Product, error)
	GetProductFileByProductID(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*model.Product, error)
	GetProducts(ctx context.Context) (model.Products, error)

	GetProductsByMerchant(ctx context.Context, merchant *model.Merchant) (model.Products, error)
	CreateProduct(ctx context.Context, req *request.ProductRequest, merchant *model.Merchant) (*model.Product, error)
}

func NewProductService(storageService StorageService, categoryService CategoryService, productRepo repo.ProductRepo,
	fileRepo repo.ProductFileRepo, txRepo repo.TransactionRepo) ProductService {
	return &productService{
		storageService:  storageService,
		categoryService: categoryService,
		productRepo:     productRepo,
		productFileRepo: fileRepo,
		txRepo:          txRepo,
	}
}

func (s *productService) GetProductById(ctx context.Context, id uint) (*model.Product, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, logger.Error(ctx, err, "error getting product by id")
	}

	return s.GetProductFileByProductID(ctx, product)
}

func (s *productService) GetProductFileByProductID(ctx context.Context, product *model.Product) (*model.Product, error) {
	var err error = nil
	product.ProductFile, err = s.productFileRepo.GetProductFileByProductID(ctx, product.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "error getting product file by product id")
	}
	return product, nil
}

func (s *productService) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	return s.productRepo.GetProductBySlug(ctx, slug)
}

func (s *productService) GetProducts(ctx context.Context) (model.Products, error) {
	return s.productRepo.GetProducts(ctx)
}

func (s *productService) GetProductsByMerchant(ctx context.Context, merchant *model.Merchant) (model.Products, error) {
	return s.productRepo.GetProductsByMerchantID(ctx, merchant.MerchantID)
}

func (s *productService) CreateProduct(ctx context.Context, req *request.ProductRequest,
	merchant *model.Merchant) (*model.Product, error) {
	product := &model.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: &req.Description,
		Slug:        slug.Make(req.Name),
		PriceAmount: req.Price,
		IsActive:    false,
		InStock:     false,
		MerchantID:  &merchant.MerchantID,
	}

	tx, err := s.txRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to begin transaction")
	}
	defer tx.End()

	if err = s.productRepo.CreateProductWithTx(ctx, tx.Transaction, product); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to create product")
	}

	productFile := &model.ProductFile{
		FileID:      uuid.New(),
		ProductID:   product.ID,
		ContentType: req.Image.ContentType(),
		Position:    0,
		IsPrimary:   true,
	}

	if err = s.productFileRepo.CreateProductFileWithTx(ctx, tx.Transaction, productFile); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to create product image")
	}

	if err = s.categoryService.AddProductCategoryWithTx(ctx, tx.Transaction, product.ID, req.CategoryIDs); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to add product category")
	}

	if err = s.storageService.UploadFile(ctx, req.Image, fmt.Sprintf("products/%s", productFile.FileID)); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to upload product image")
	}

	product.ProductFile = productFile
	tx.CanCommit = true
	return product, nil
}
