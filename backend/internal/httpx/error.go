package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"alwis.dev/selectify/internal/logger"
)

var (
	ErrInvalidRequest = fmt.Errorf("invalid_request")
	ErrInvalidCountry = errors.New("invalid_country")

	ErrInternalServer = fmt.Errorf("internal_server_error")
	ErrUserNotFound   = errors.New("user_not _found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrBadRequest     = errors.New("bad_request")
)

var (
	ErrUserAlreadyExists = fmt.Errorf("USER_ALREADY_EXISTS")
)

// SendError return http.StatusBadRequest with an error message
func SendError(w http.ResponseWriter, err error) {
	if nErr := SendJson(w, http.StatusBadRequest, &StatusResponse{Status: "error", Message: err.Error()}); nErr != nil {
		_ = logger.Error(context.Background(), nErr, "Failed to send json response")
	}
}

func SendUnAuthorized(w http.ResponseWriter) {
	DeleteSessionCookies(w)
	if err := SendJson(w, http.StatusUnauthorized, &StatusResponse{Status: "error", Message: ErrUnauthorized.Error()}); err != nil {
		_ = logger.Error(context.Background(), err, "Failed to send json response")
	}
}

func SendForbidden(w http.ResponseWriter) {
	if err := SendJson(w, http.StatusForbidden, &StatusResponse{Status: "error", Message: ErrForbidden.Error()}); err != nil {
		_ = logger.Error(context.Background(), err, "Failed to send json response")
	}
}

func SendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func SendUnAuthorizedWithNotFoud(w http.ResponseWriter) {
	DeleteSessionCookies(w)
	SendNotFound(w)
}
