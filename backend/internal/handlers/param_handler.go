package handlers

import (
	"net/http"

	"alwis.dev/selectify/internal/model"
)

type ParamHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *model.LoggedInSession)
}
