package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	CountryRepo     repo.CountryRepo
	MerchantRepo    repo.MerchantRepo
	UserRepo        repo.UserRepo
	UserDeviceRepo  repo.UserDeviceRepo
	UserSessionRepo repo.UserSessionRepo
	UserRoleRepo    repo.UserRoleRepo
	TxRepo          repo.TransactionRepo
}

func NewRepository() *Repo {
	r := new(Repo)

	r.CountryRepo = repo.NewCountryRepo(appEnv.dbConn)
	r.MerchantRepo = repo.NewMerchantRepo(appEnv.dbConn)
	r.UserRepo = repo.NewUserRepo(appEnv.dbConn)
	r.UserDeviceRepo = repo.NewUserDeviceRepo(appEnv.dbConn)
	r.UserRoleRepo = repo.NewUserRoleRepo(appEnv.dbConn)
	r.UserSessionRepo = repo.NewUserSessionRepo(appEnv.dbConn)
	r.TxRepo = repo.NewTransactionRepository(appEnv.dbConn)
	return r
}
