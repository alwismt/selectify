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
		nil,
		appEnv.repo.TxRepo,
		appEnv.repo.UserFileRepo,
		appEnv.repo.UserRepo,
		service.NewEmailService(),
		service.NewGeoIPService(),
	)
	return svc
}
