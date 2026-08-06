package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	CartRepo            repo.CartRepo
	OrderRepo           repo.OrderRepo
	ProductRepo         repo.ProductRepo
	ProductFileRepo     repo.ProductFileRepo
	ProductVariantsRepo repo.ProductVariantsRepo
	PaymentRepo         repo.PaymentRepo
	TxRepo              repo.TransactionRepo
	UserRepo            repo.UserRepo
	UserAddressRepo     repo.UserAddressRepo
	UserFileRepo        repo.UserFileRepo
	UserRoleRepo        repo.UserRoleRepo
	UserSessionRepo     repo.UserSessionRepo
}

func NewRepository() *Repo {
	repository := new(Repo)

	repository.CartRepo = repo.NewCartRepo(appEnv.dbConn)
	repository.OrderRepo = repo.NewOrderRepo(appEnv.dbConn)
	repository.ProductRepo = repo.NewProductRepo(appEnv.dbConn)
	repository.ProductFileRepo = repo.NewProductFileRepo(appEnv.dbConn)
	repository.ProductVariantsRepo = repo.NewProductVariantRepo(appEnv.dbConn)
	repository.TxRepo = repo.NewTransactionRepository(appEnv.dbConn)
	repository.UserRepo = repo.NewUserRepo(appEnv.dbConn)
	repository.UserFileRepo = repo.NewUserFileRepo(appEnv.dbConn)
	repository.UserAddressRepo = repo.NewUserAddressRepo(appEnv.dbConn)
	repository.UserRoleRepo = repo.NewUserRoleRepo(appEnv.dbConn)
	repository.UserSessionRepo = repo.NewUserSessionRepo(appEnv.dbConn)
	repository.PaymentRepo = repo.NewPaymentRepo(appEnv.dbConn)

	return repository
}
