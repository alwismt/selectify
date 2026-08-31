package handlers

import (
	"net/http"

	selectifyHandlers "alwis.dev/selectify/internal/handlers"
	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/permission"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/merchant-pkg/app"
)

type MerchantAccessHandlerFunc func(http.ResponseWriter, *http.Request, *model.LoggedInSession, *model.Merchant)

type merchantAccessHandler struct {
	requiredPermissions []permission.Permission
	merchantRepo        repo.MerchantRepo

	handlerFunc MerchantAccessHandlerFunc
}

func MerchantAccessHandler(fn MerchantAccessHandlerFunc, requiredPermissions ...permission.Permission) selectifyHandlers.ParamHandler {
	return merchantAccessHandler{}.init(fn, requiredPermissions...)
}

func (h merchantAccessHandler) init(fn MerchantAccessHandlerFunc, requiredPermissions ...permission.Permission) merchantAccessHandler {
	h.requiredPermissions = requiredPermissions
	h.merchantRepo = app.Repository().MerchantRepo
	h.handlerFunc = fn
	return h
}

func (h merchantAccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	if s.UserRole.MerchantRole == nil || !permission.HasMerchantPermission(*s.UserRole.MerchantRole, h.requiredPermissions...) {
		_ = logger.Errorf(r.Context(), nil, "Merchant doesn't have permission. Required permission: %s, User Id ID: %d",
			h.requiredPermissions, s.UserId)
		httpx.SendNotFound(w)
		return
	}

	merchant, err := h.merchantRepo.GetMerchant(r.Context(), *s.UserRole.MerchantID)
	if err != nil {
		_ = logger.Errorf(r.Context(), err, "Failed to get merchant %d", *s.UserRole.MerchantID)
		httpx.SendNotFound(w)
		return
	}

	h.handlerFunc(w, r, s, merchant)
}
