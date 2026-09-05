package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	CountryRepo         repo.CountryRepo
	CategoryRepo        repo.CategoryRepo
	MerchantRepo        repo.MerchantRepo
	ProductRepo         repo.ProductRepo
	ProductFileRepo     repo.ProductFileRepo
	ProductVariantsRepo repo.ProductVariantsRepo
	TxRepo              repo.TransactionRepo
	UserRepo            repo.UserRepo
	UserDeviceRepo      repo.UserDeviceRepo
	UserSessionRepo     repo.UserSessionRepo
	UserRoleRepo        repo.UserRoleRepo
}

func NewRepository() *Repo {
	r := new(Repo)

	r.CountryRepo = repo.NewCountryRepo(appEnv.dbConn)
	r.CategoryRepo = repo.NewCategoryRepo(appEnv.dbConn)
	r.MerchantRepo = repo.NewMerchantRepo(appEnv.dbConn)
	r.ProductRepo = repo.NewProductRepo(appEnv.dbConn)
	r.ProductVariantsRepo = repo.NewProductVariantRepo(appEnv.dbConn)
	r.ProductFileRepo = repo.NewProductFileRepo(appEnv.dbConn)
	r.TxRepo = repo.NewTransactionRepository(appEnv.dbConn)
	r.UserRepo = repo.NewUserRepo(appEnv.dbConn)
	r.UserDeviceRepo = repo.NewUserDeviceRepo(appEnv.dbConn)
	r.UserRoleRepo = repo.NewUserRoleRepo(appEnv.dbConn)
	r.UserSessionRepo = repo.NewUserSessionRepo(appEnv.dbConn)
	return r
}
