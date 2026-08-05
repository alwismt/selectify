package middleware

import (
	"net/http"
)

func HeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			// Disable caching
			headers.Set("Cache-Control", "no-cache, no-store, must-revalidate")

			headers.Set("X-Frame-Options", "SAMEORIGIN")

			headers.Set("X-Content-Type-Options", "nosniff")

			// Strict-Transport-Security
			//headers.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

			headers.Set("Content-Security-Policy", "default-src 'self'; object-src 'none';")

			headers.Set("X-Robots-Tag", "noindex, nofollow")

			headers.Set("Access-Control-Allow-Origin", "https://*.alwis.dev")
			headers.Set("Vary", "Origin")
			headers.Set("Access-Control-Allow-Credentials", "true")
			headers.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Token")

			next.ServeHTTP(w, r)
		})
	}
}
