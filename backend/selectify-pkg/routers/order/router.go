package order

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/params"
	"alwis.dev/selectify/selectify-pkg/handlers"
)

func Register(r chi.Router) {
	c := new(controller).init()
	r.Route("/orders", func(r chi.Router) {
		// route POST /api/v1/orders
		r.Method(http.MethodPost, "/", handlers.UserSessionHandlerFunc(c.CreateOrder))

		// route GET /api/v1/orders
		r.Method(http.MethodGet, "/", handlers.UserSessionHandlerFunc(c.GetOrders))

		// route PUT /api/v1/orders/{order_id}/address
		r.Method(http.MethodPut, "/"+params.OrderIdParam+"/address", handlers.UserSessionHandlerFunc(c.SetOrderShippingAddress))
	})
}
