package app

import "alwis.dev/selectify/internal/repo"

type Repo struct {
	MerchantRepo    repo.MerchantRepo
	UserRepo        repo.UserRepo
	UserDeviceRepo  repo.UserDeviceRepo
	UserSessionRepo repo.UserSessionRepo
	UserRoleRepo    repo.UserRoleRepo
}

func NewRepository() *Repo {
	r := new(Repo)

	r.MerchantRepo = repo.NewMerchantRepo(appEnv.dbConn)
	r.UserRepo = repo.NewUserRepo(appEnv.dbConn)
	r.UserDeviceRepo = repo.NewUserDeviceRepo(appEnv.dbConn)
	r.UserSessionRepo = repo.NewUserSessionRepo(appEnv.dbConn)
	r.UserRoleRepo = repo.NewUserRoleRepo(appEnv.dbConn)
	return r
}
