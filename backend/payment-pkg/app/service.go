package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	PaymentService service.PaymentService
}

func NewService() *SVC {
	svc := new(SVC)
	svc.PaymentService = service.NewPaymentService(appEnv.repo.OrderRepo, appEnv.repo.PaymentRepo)
	return svc
}
