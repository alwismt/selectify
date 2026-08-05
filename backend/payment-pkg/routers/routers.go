package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/middleware"
)

func CreateHandler() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/pay", func(r chi.Router) {
		registerRoutes(r)
	})

	return registerHandlers(r)
}

func registerRoutes(r chi.Router) {
	c := new(controller).init()
	r.Route("/webhooks", func(r chi.Router) {
		r.Post("/stripe", c.IntentUpdate)
	})
}

func registerHandlers(h http.Handler) http.Handler {
	corsMw := middleware.CORSMiddleware()
	headersMw := middleware.HeadersMiddleware()
	return headersMw(corsMw(h))
}
