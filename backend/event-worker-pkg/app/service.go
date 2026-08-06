package app

import (
	"alwis.dev/selectify/internal/service"
)

type SVC struct {
	UserService service.UserService
}

func NewService() *SVC {
	svc := new(SVC)
	svc.UserService = service.NewUserService(
		noopStorage{},
		appEnv.repo.TxRepo,
		appEnv.repo.UserFileRepo,
		appEnv.repo.UserRepo,
	)
	return svc
}
