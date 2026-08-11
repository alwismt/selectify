package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/middleware"
	"alwis.dev/selectify/merchant-pkg/routers/merchant"
	"alwis.dev/selectify/merchant-pkg/routers/product"
)

func CreateHandler() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/merchant", func(r chi.Router) {
		registerRoutes(r)
	})
	return registerHandlers(r)
}

func registerRoutes(r chi.Router) {
	merchant.RegisterRoutes(r)
	product.RegisterRoutes(r)
}

func registerHandlers(h http.Handler) http.Handler {
	corsMw := middleware.CORSMiddleware()
	headersMw := middleware.HeadersMiddleware()
	return headersMw(corsMw(h))
}
