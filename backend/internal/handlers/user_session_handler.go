package handlers

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
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
	userRepo := app.Repository().UserRepo

	ctx := r.Context()

	c, err := r.Cookie(httpx.Cookie1)
	if err != nil || c == nil || c.Value == "" {
		logger.Warnf(ctx, "session cookie not found or invalid, %v", err)
		httpx.SendUnAuthorized(w)
		return
	}

	s, err := sessionRepo.GetBySessionId(ctx, c.Value)
	if err != nil {
		logger.Warn(ctx, "session not found")
		httpx.SendUnAuthorized(w)
		return
	}

	//if s.UserAgent != r.UserAgent() {
	//	fmt.Println()
	//	fmt.Println("1 user agent is", s.UserAgent)
	//	fmt.Println("2 user agent is", r.UserAgent())
	//	logger.Warn(ctx, "user agent not match")
	//	if err = sessionRepo.RevokeSession(ctx, s.SessionId); err != nil {
	//		logger.Warn(ctx, "failed to revoke session")
	//	}
	//	httpx.DeleteSessionCookies(w)
	//	httpx.SendUnAuthorized(w)
	//	return
	//}

	if s.IsExpired() {
		logger.Warn(ctx, "session is expired")
		httpx.SendUnAuthorized(w)
		return
	}

	s.User, err = userRepo.GetUserById(ctx, s.UserId)
	if err != nil {
		logger.Warn(ctx, "user not found")
		httpx.SendUnAuthorized(w)
	}

	h.handler.ServeHTTP(w, r, s)
}
