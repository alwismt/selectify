package handlers

import (
	"net/http"

	"alwis.dev/selectify/internal/handlers"
)

func UserSessionHandler(h handlers.ParamHandler) http.Handler {
	return handlers.UserSessionHandler(h)
}
