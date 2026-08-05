package app

import "alwis.dev/selectify/internal/service"

type SVC struct {
	AuthService            service.AuthService
	CartService            service.CartService
	JWTService             service.JWTService
	OrderService           service.OrderService
	ProductService         service.ProductService
	ProductVariantsService service.ProductVariantsService
	PaymentService         service.PaymentService
}

func NewService() *SVC {
	svc := new(SVC)
	svc.JWTService = service.NewJWTService()
	svc.AuthService = service.NewAuthService(svc.JWTService, appEnv.repo.UserRepo, appEnv.repo.TxRepo, appEnv.repo.UserRoleRepo,
		appEnv.repo.UserSessionRepo)
	svc.ProductService = service.NewProductService(appEnv.repo.ProductRepo, appEnv.repo.ProductFileRepo)
	svc.ProductVariantsService = service.NewProductVariantsService(appEnv.repo.ProductVariantsRepo, appEnv.repo.ProductFileRepo)
	svc.CartService = service.NewCartService(appEnv.repo.CartRepo, appEnv.repo.ProductRepo, appEnv.repo.ProductVariantsRepo)
	svc.OrderService = service.NewOrderService(svc.CartService, appEnv.repo.OrderRepo, appEnv.repo.TxRepo, appEnv.repo.ProductVariantsRepo,
		appEnv.grpcCli.PaymentClient, appEnv.repo.UserAddressRepo)
	svc.PaymentService = service.NewPaymentService(appEnv.repo.OrderRepo, appEnv.repo.PaymentRepo)
	return svc
}
