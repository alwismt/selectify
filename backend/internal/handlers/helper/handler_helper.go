package handlers_helper

import "alwis.dev/selectify/internal/repo"

type SessionHandlerHelper struct {
	UserSessionRepo repo.UserSessionRepo
	DeviceRepo      repo.UserDeviceRepo
	UserRepo        repo.UserRepo
	UserRoleRepo    repo.UserRoleRepo
}

func NewSessionHandlerHelper(uSessionRepo repo.UserSessionRepo, deviceRepo repo.UserDeviceRepo, userRepo repo.UserRepo,
	userRoleRepo repo.UserRoleRepo) *SessionHandlerHelper {
	return &SessionHandlerHelper{
		UserSessionRepo: uSessionRepo,
		DeviceRepo:      deviceRepo,
		UserRepo:        userRepo,
		UserRoleRepo:    userRoleRepo,
	}
}
