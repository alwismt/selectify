package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	CountryService         service.CountryService
	CategoryService        service.CategoryService
	MerchantService        service.MerchantService
	StorageService         service.StorageService
	ProductService         service.ProductService
	ProductFileService     service.ProductFileService
	ProductVariantsService service.ProductVariantsService
}

func NewService() *SVC {
	svc := new(SVC)

	svc.StorageService = service.NewStorageService()
	svc.CountryService = service.NewCountryService(appEnv.repo.CountryRepo)
	svc.CategoryService = service.NewCategoryService(appEnv.repo.CategoryRepo)
	svc.MerchantService = service.NewMerchantService(svc.StorageService, appEnv.repo.CountryRepo,
		appEnv.repo.MerchantRepo, appEnv.repo.TxRepo)
	svc.ProductService = service.NewProductService(svc.StorageService, svc.CategoryService, appEnv.repo.ProductRepo, appEnv.repo.ProductFileRepo, appEnv.repo.TxRepo)
	svc.ProductFileService = service.NewProductFileService(appEnv.repo.ProductFileRepo)
	svc.ProductVariantsService = service.NewProductVariantsService(appEnv.repo.ProductVariantsRepo, appEnv.repo.ProductFileRepo)
	return svc
}
