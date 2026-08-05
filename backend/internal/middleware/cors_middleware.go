package middleware

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
)

const allowedOrigin = "https://selectify.alwis.dev"

func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" {
				if !isAllowedOrigin(origin) {
					http.Error(w, "CORS: Origin not allowed", http.StatusForbidden)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.Header().Set("Vary", "Origin")
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				_ = httpx.SendJson(w, http.StatusForbidden, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isAllowedOrigin(origin string) bool {
	return origin == allowedOrigin
}
