package handlers

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/security"
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type sessionHandler struct {
	handler ParamHandler
}

type ParamHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *model.UserSession)
}

func UserSessionHandlerFunc(fn func(http.ResponseWriter, *http.Request, *model.UserSession)) http.HandlerFunc {
	h := initSessionHandler(funcHandler{fn: fn})
	return h.ServeHTTP
}

type funcHandler struct {
	fn func(http.ResponseWriter, *http.Request, *model.UserSession)
}

func (h funcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	h.fn(w, r, s)
}

func UserSessionHandler(h ParamHandler) http.Handler {
	return initSessionHandler(h)
}

func initSessionHandler(h ParamHandler) http.Handler {
	return &sessionHandler{
		handler: h,
	}
}

func (h sessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionRepo := app.Repository().UserSessionRepo
	deviceRepo := app.Repository().UserDeviceRepo
	userRepo := app.Repository().UserRepo
	userRoleRepo := app.Repository().UserRoleRepo

	ctx := r.Context()

	c, err := r.Cookie(httpx.Cookie1)
	if err != nil || c == nil || c.Value == "" {
		logger.Warnf(ctx, "session cookie not found or invalid, %v", err)
		httpx.SendUnAuthorized(w)
		return
	}

	s, err := sessionRepo.GetByTokenHash(ctx, security.HashToken(c.Value))
	if err != nil {
		logger.Warn(ctx, "session not found")
		httpx.SendUnAuthorized(w)
		return
	}

	if s.IsExpired() {
		logger.Warn(ctx, "session is expired")
		httpx.SendUnAuthorized(w)
		return
	}

	if renewed, renewErr := sessionRepo.RenewIfStale(ctx, s.SessionId, service.SessionIdleTTL, service.ActivityThrottle); renewErr != nil {
		logger.Warn(ctx, "failed to renew session")
	} else if renewed != nil {
		s = renewed
	}

	if s.DeviceId != nil {
		_ = deviceRepo.TouchIfStale(ctx, *s.DeviceId, r.UserAgent(), httpx.ClientIP(r), service.ActivityThrottle)
	}

	s.User, err = userRepo.GetUserById(ctx, s.UserId)
	if err != nil || s.User == nil {
		logger.Warn(ctx, "user not found")
		httpx.SendUnAuthorized(w)
		return
	}

	s.UserRole, err = userRoleRepo.GetUserRoleByUserID(ctx, s.UserId)
	if err != nil || s.UserRole == nil {
		logger.Warn(ctx, "user role not found")
		httpx.SendUnAuthorized(w)
		return
	}
	if !s.UserRole.Role.Valid() {
		logger.Warn(ctx, "user role not valid")
		httpx.SendUnAuthorized(w)
		return
	}

	h.handler.ServeHTTP(w, r, s)
}
