package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	MerchantService service.MerchantService
}

func NewService() *SVC {
	svc := new(SVC)

	svc.MerchantService = service.NewMerchantService(appEnv.repo.MerchantRepo)

	return svc
}
