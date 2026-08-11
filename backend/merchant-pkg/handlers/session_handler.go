package handlers

import (
	"net/http"

	"alwis.dev/selectify/internal/handlers"
	handlersHelper "alwis.dev/selectify/internal/handlers/helper"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/merchant-pkg/app"
)

func MerchantSessionHandler(h handlers.ParamHandler) http.Handler {
	sHandlerHelper := handlersHelper.NewSessionHandlerHelper(app.Repository().UserSessionRepo, app.Repository().UserDeviceRepo,
		app.Repository().UserRepo, app.Repository().UserRoleRepo)
	return handlers.MerchantSessionHandler(h, sHandlerHelper)
}

func MerchantSessionHandlerFunc(fn func(http.ResponseWriter, *http.Request, *model.LoggedInSession)) http.HandlerFunc {
	sHandlerHelper := handlersHelper.NewSessionHandlerHelper(app.Repository().UserSessionRepo, app.Repository().UserDeviceRepo,
		app.Repository().UserRepo, app.Repository().UserRoleRepo)
	return handlers.MerchantSessionHandlerFunc(fn, sHandlerHelper)
}
