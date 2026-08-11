package handlers

import (
	"net/http"

	handlersHelper "alwis.dev/selectify/internal/handlers/helper"
	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/security"
	"alwis.dev/selectify/internal/service"
)

type sessionHandler struct {
	mustMerchant bool
	mustAdmin    bool
	handler      ParamHandler

	deviceRepo   repo.UserDeviceRepo
	sessionRepo  repo.UserSessionRepo
	userRepo     repo.UserRepo
	userRoleRepo repo.UserRoleRepo
}

type ParamHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *model.LoggedInSession)
}

func UserSessionHandlerFunc(fn func(http.ResponseWriter, *http.Request, *model.LoggedInSession), helper *handlersHelper.SessionHandlerHelper) http.HandlerFunc {
	h := initSessionHandler(funcHandler{fn: fn}, false, false, helper)
	return h.ServeHTTP
}

func UserSessionHandler(h ParamHandler, helper *handlersHelper.SessionHandlerHelper) http.Handler {
	return initSessionHandler(h, false, false, helper)
}

func MerchantSessionHandler(h ParamHandler, helper *handlersHelper.SessionHandlerHelper) http.Handler {
	return initSessionHandler(h, true, false, helper)
}

func MerchantSessionHandlerFunc(fn func(http.ResponseWriter, *http.Request, *model.LoggedInSession), helper *handlersHelper.SessionHandlerHelper) http.HandlerFunc {
	h := initSessionHandler(funcHandler{fn: fn}, true, false, helper)
	return h.ServeHTTP
}

type funcHandler struct {
	fn func(http.ResponseWriter, *http.Request, *model.LoggedInSession)
}

func (h funcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	h.fn(w, r, s)
}

func initSessionHandler(h ParamHandler, mustMerchant, mustAdmin bool, helper *handlersHelper.SessionHandlerHelper) http.Handler {
	return &sessionHandler{
		handler:      h,
		mustMerchant: mustMerchant,
		mustAdmin:    mustAdmin,

		deviceRepo:   helper.DeviceRepo,
		sessionRepo:  helper.UserSessionRepo,
		userRepo:     helper.UserRepo,
		userRoleRepo: helper.UserRoleRepo,
	}
}

func (h sessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, err := h.getUserAndRoleBySession(w, r)
	if err != nil {
		return
	}
	if s == nil || !s.UserRole.Role.Valid() {
		logger.Warn(r.Context(), "User role not valid")
		httpx.SendUnAuthorized(w)
		return
	}

	if h.mustMerchant {
		if !s.UserRole.Role.IsMerchant() {
			logger.Warn(r.Context(), "User is not a merchant")
			httpx.SendUnAuthorizedWithNotFoud(w)
			return
		}
	}

	h.handler.ServeHTTP(w, r, s)
}

func (h sessionHandler) getUserAndRoleBySession(w http.ResponseWriter, r *http.Request) (*model.LoggedInSession, error) {

	ctx := r.Context()

	c, err := r.Cookie(httpx.Cookie1)
	if err != nil || c == nil || c.Value == "" {
		logger.Warnf(ctx, "session cookie not found or invalid, %v", err)
		httpx.SendUnAuthorized(w)
		return nil, err
	}

	s, err := h.sessionRepo.GetByTokenHash(ctx, security.HashToken(c.Value))
	if err != nil {
		logger.Warn(ctx, "session not found")
		httpx.SendUnAuthorized(w)
		return nil, err
	}

	if s.IsExpired() {
		logger.Warn(ctx, "session is expired")
		httpx.SendUnAuthorized(w)
		return nil, err
	}

	if renewed, renewErr := h.sessionRepo.RenewIfStale(ctx, s.SessionId, service.SessionIdleTTL, service.ActivityThrottle); renewErr != nil {
		logger.Warn(ctx, "failed to renew session")
	} else if renewed != nil {
		s = renewed
	}

	if s.DeviceId != nil {
		_ = h.deviceRepo.TouchIfStale(ctx, *s.DeviceId, r.UserAgent(), httpx.ClientIP(r), service.ActivityThrottle)
	}

	s.User, err = h.userRepo.GetUserById(ctx, s.UserId)
	if err != nil || s.User == nil {
		logger.Warn(ctx, "user not found")
		httpx.SendUnAuthorized(w)
		return nil, err
	}

	s.UserRole, err = h.userRoleRepo.GetUserRoleByUserID(ctx, s.UserId)
	if err != nil || s.UserRole == nil {
		logger.Warn(ctx, "user role not found")
		httpx.SendUnAuthorized(w)
		return nil, err
	}

	return s, nil
}
