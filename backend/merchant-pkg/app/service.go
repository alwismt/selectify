package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	CountryService  service.CountryService
	MerchantService service.MerchantService
	StorageService  service.StorageService
}

func NewService() *SVC {
	svc := new(SVC)

	svc.StorageService = service.NewStorageService()
	svc.CountryService = service.NewCountryService(appEnv.repo.CountryRepo)
	svc.MerchantService = service.NewMerchantService(svc.StorageService, appEnv.repo.CountryRepo,
		appEnv.repo.MerchantRepo, appEnv.repo.TxRepo)

	return svc
}
