package auth

import (
	"errors"
	"fmt"
	"net/http"

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
	clientIP := httpx.ClientIP(r)
	deviceToken := deviceTokenFromRequest(r)

	s, err := c.authService.RegisterUser(ctx, req, userAgent, clientIP, deviceToken)
	if err != nil {
		_ = logger.Errorf(ctx, err, "Failed to register user")

		if errors.Is(err, httpx.ErrUserAlreadyExists) {
			httpx.SendError(w, httpx.ErrUserAlreadyExists)
			return
		}
		httpx.SendError(w, httpx.ErrInvalidRequest)
		return
	}

	httpx.SetSessionCookies(s, w)

	httpx.SendStatus(w, http.StatusCreated, "User logged in")
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

	s, err := c.authService.LoginUser(ctx, req, r.UserAgent(), httpx.ClientIP(r), deviceTokenFromRequest(r))
	if err != nil {
		_ = logger.Errorf(ctx, err, "Failed to login")
		if errors.Is(err, httpx.ErrUserNotFound) {
			httpx.SendError(w, httpx.ErrUserNotFound)
			return
		}
		if errors.Is(err, service.ErrUserDeactivated) {
			httpx.SendError(w, service.ErrUserDeactivated)
			return
		}
		httpx.SendStatus(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httpx.SetSessionCookies(s, w)

	httpx.SendStatus(w, http.StatusCreated, "User logged in")
}

func deviceTokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(httpx.DeviceCookie)
	if err != nil || c == nil {
		return ""
	}
	return c.Value
}

func (c *controller) Logout(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	if err := c.authService.UserLogout(r.Context(), s.SessionId); err != nil {
		httpx.SendError(w, fmt.Errorf("failed logout"))
		return
	}

	httpx.DeleteSessionCookies(w)
	httpx.SendStatus(w, http.StatusOK, "User logged out")
	return
}

func (c *controller) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(request.ForgotPasswordRequest)

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

	if err := c.authService.ForgetPassword(ctx, req.Email, httpx.ClientIP(r), r.UserAgent()); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to process forgot password")
	}

	httpx.SendStatus(w, http.StatusOK, "If an account exists with this email address, a password reset link has been sent.")
}

func (c *controller) ValidateResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		httpx.SendError(w, service.ErrInvalidResetToken)
		return
	}

	if err := c.authService.ValidateResetToken(ctx, token); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to validate reset password token")
		if errors.Is(err, service.ErrInvalidResetToken) {
			httpx.SendError(w, service.ErrInvalidResetToken)
			return
		}
		httpx.SendStatus(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httpx.SendStatus(w, http.StatusOK, "Token is valid")
}

func (c *controller) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(request.ResetPasswordRequest)

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

	if err := c.authService.ResetPassword(ctx, req.Token, req.Password, httpx.ClientIP(r), r.UserAgent()); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to reset password")
		if errors.Is(err, service.ErrInvalidResetToken) {
			httpx.SendError(w, service.ErrInvalidResetToken)
			return
		}
		httpx.SendStatus(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	httpx.SendStatus(w, http.StatusOK, "Password has been reset successfully.")
}
