package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"alwis.dev/selectify/internal/logger"
)

var (
	ErrInvalidRequest = fmt.Errorf("INVALID_REQUEST")
	ErrInternalServer = fmt.Errorf("INTERNAL_SERVER")
	ErrUserNotFound   = errors.New("user not found")
	ErrUnauthorized   = errors.New("unauthorized")
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

func SendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func SendUnAuthorizedWithNotFoud(w http.ResponseWriter) {
	DeleteSessionCookies(w)
	SendNotFound(w)
}
