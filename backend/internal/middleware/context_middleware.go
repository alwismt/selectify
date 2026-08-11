package middleware

import (
	"net/http"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/security"
)

type requestContextKey uint

const (
	sessionKey requestContextKey = iota + 4783478488
)

type sessionContext struct {
	*model.LoggedInSession
	Valid bool
}

// ContextMiddleware is unused by the main router; session auth uses UserSessionHandler.
// Kept for optional future context attachment.
func ContextMiddleware(sessionRepo repo.UserSessionRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			c, err := r.Cookie("slf")
			if err == nil && c.Value != "" {
				_, _ = sessionRepo.GetByTokenHash(ctx, security.HashToken(c.Value))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
