package app

import (
	"alwis.dev/selectify/internal/service"
)

type SVC struct {
	AuthService            service.AuthService
	CartService            service.CartService
	CategoryService        service.CategoryService
	EventPublisher         service.EventPublisher
	OrderService           service.OrderService
	ProductService         service.ProductService
	ProductVariantsService service.ProductVariantsService
	PaymentService         service.PaymentService
	StorageService         service.StorageService
	UserService            service.UserService
	SiteConfigService      service.SiteConfigService
}

func NewService() *SVC {
	svc := new(SVC)
	svc.EventPublisher = service.NewSQSEventPublisher(appEnv.repo.EventRepo)
	svc.AuthService = service.NewAuthService(svc.EventPublisher, appEnv.repo.UserRepo, appEnv.repo.TxRepo,
		appEnv.repo.UserRoleRepo, appEnv.repo.UserSessionRepo, appEnv.repo.UserDeviceRepo, appEnv.repo.PasswordResetRepo)
	svc.CategoryService = service.NewCategoryService(appEnv.repo.CategoryRepo)
	svc.StorageService = service.NewStorageService()
	svc.ProductService = service.NewProductService(svc.StorageService, svc.CategoryService, appEnv.repo.ProductRepo,
		appEnv.repo.ProductFileRepo, appEnv.repo.TxRepo)
	svc.ProductVariantsService = service.NewProductVariantsService(appEnv.repo.ProductVariantsRepo, appEnv.repo.ProductFileRepo)
	svc.CartService = service.NewCartService(appEnv.repo.CartRepo, appEnv.repo.ProductRepo, appEnv.repo.ProductVariantsRepo)
	svc.OrderService = service.NewOrderService(svc.CartService, appEnv.repo.OrderRepo, appEnv.repo.TxRepo, appEnv.repo.ProductVariantsRepo,
		appEnv.grpcCli.PaymentClient, appEnv.repo.UserAddressRepo, appEnv.repo.CurrencyRepo)
	svc.PaymentService = service.NewPaymentService(appEnv.repo.OrderRepo, appEnv.repo.PaymentRepo)
	svc.UserService = service.NewUserService(service.NewStorageService(), appEnv.repo.TxRepo, appEnv.repo.UserFileRepo,
		appEnv.repo.UserRepo, nil, nil)
	svc.SiteConfigService = service.NewSiteConfigService(appEnv.repo.CurrencyRepo)
	return svc
}
