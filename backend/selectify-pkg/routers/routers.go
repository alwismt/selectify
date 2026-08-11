package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/middleware"
	"alwis.dev/selectify/selectify-pkg/routers/auth"
	"alwis.dev/selectify/selectify-pkg/routers/cart"
	"alwis.dev/selectify/selectify-pkg/routers/order"
	"alwis.dev/selectify/selectify-pkg/routers/product"
	"alwis.dev/selectify/selectify-pkg/routers/status"
	"alwis.dev/selectify/selectify-pkg/routers/user"
)

func CreateHandler() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		registerRoutes(r)
		//registerRoutesWithCtx(r)
	})

	return registerHandlers(r)
}

func registerRoutes(r chi.Router) {
	status.Register(r)
	product.Register(r)
	auth.Register(r)
	user.Register(r)
	cart.Register(r)
	order.Register(r)
}

func registerHandlers(h http.Handler) http.Handler {
	corsMw := middleware.CORSMiddleware()
	headersMw := middleware.HeadersMiddleware()
	return headersMw(corsMw(h))
}
