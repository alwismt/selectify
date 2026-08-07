package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	AuthService            service.AuthService
	CartService            service.CartService
	EventPublisher         service.EventPublisher
	OrderService           service.OrderService
	ProductService         service.ProductService
	ProductVariantsService service.ProductVariantsService
	PaymentService         service.PaymentService
	UserService            service.UserService
}

func NewService() *SVC {
	svc := new(SVC)
	svc.EventPublisher = service.NewSQSEventPublisher(appEnv.repo.EventRepo)
	svc.AuthService = service.NewAuthService(svc.EventPublisher, appEnv.repo.UserRepo, appEnv.repo.TxRepo,
		appEnv.repo.UserRoleRepo, appEnv.repo.UserSessionRepo, appEnv.repo.UserDeviceRepo, appEnv.repo.PasswordResetRepo)
	svc.ProductService = service.NewProductService(appEnv.repo.ProductRepo, appEnv.repo.ProductFileRepo)
	svc.ProductVariantsService = service.NewProductVariantsService(appEnv.repo.ProductVariantsRepo, appEnv.repo.ProductFileRepo)
	svc.CartService = service.NewCartService(appEnv.repo.CartRepo, appEnv.repo.ProductRepo, appEnv.repo.ProductVariantsRepo)
	svc.OrderService = service.NewOrderService(svc.CartService, appEnv.repo.OrderRepo, appEnv.repo.TxRepo, appEnv.repo.ProductVariantsRepo,
		appEnv.grpcCli.PaymentClient, appEnv.repo.UserAddressRepo)
	svc.PaymentService = service.NewPaymentService(appEnv.repo.OrderRepo, appEnv.repo.PaymentRepo)
	svc.UserService = service.NewUserService(service.NewStorageService(), appEnv.repo.TxRepo, appEnv.repo.UserFileRepo, appEnv.repo.UserRepo, nil, nil)
	return svc
}
