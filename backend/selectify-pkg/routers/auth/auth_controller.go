package auth

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/service"
)

func (c *controller) UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(request.UserRegisterRequest)

	if err := httpx.MustDecodeJson(w, r, req); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to decode body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		fieldErrors := httpx.ValidationErrorsToMap(err)
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	if err := req.Validate(); err != nil {
		_ = logger.Errorf(ctx, err, "Password validation error")
		httpx.SendError(w, err)
		return
	}
	userAgent := r.UserAgent()
	clientIP := c.getClientIP(r)

	s, err := c.authService.RegisterUser(ctx, req, userAgent, clientIP)
	if err != nil {
		_ = logger.Errorf(ctx, err, "Failed to register user")

		if errors.Is(err, httpx.ErrUserAlreadyExists) {
			httpx.SendError(w, httpx.ErrUserAlreadyExists)
			return
		}
		httpx.SendError(w, httpx.ErrInvalidRequest)
	}

	httpx.SetSessionCookies(s, w)

	httpx.SendStatus(w, http.StatusCreated, "User logged in")
	return
}

func (c *controller) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := new(request.LoginRequest)

	if err := httpx.MustDecodeJson(w, r, req); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to decode body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		fieldErrors := httpx.ValidationErrorsToMap(err)
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	s, err := c.authService.LoginUser(ctx, req, r.UserAgent(), c.getClientIP(r))
	if err != nil {
		_ = logger.Errorf(ctx, err, "Failed to login")
		if errors.Is(err, httpx.ErrUserNotFound) {
			httpx.SendError(w, httpx.ErrUserNotFound)
			return
		}
		if errors.Is(err, httpx.ErrInternalServer) {
			httpx.SendError(w, service.ErrUserDeactivated)
			return
		}
		httpx.SendStatus(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httpx.SetSessionCookies(s, w)

	httpx.SendStatus(w, http.StatusCreated, "User logged in")
	return
}

func (c *controller) getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (c *controller) Logout(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	if err := c.authService.UserLogout(r.Context(), s.SessionId); err != nil {
		httpx.SendError(w, fmt.Errorf("failed logout"))
		return
	}

	httpx.DeleteSessionCookies(w)
	httpx.SendStatus(w, http.StatusOK, "User logged out")
	return
}
