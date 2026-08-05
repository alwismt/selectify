package cart

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/handlers"
	"alwis.dev/selectify/internal/params"
	sfHandlers "alwis.dev/selectify/selectify-pkg/handlers"
)

func Register(r chi.Router) {
	c := new(controller).init()
	r.Route("/cart", func(r chi.Router) {
		// route GET /api/v1/cart
		r.Method(http.MethodGet, "/", handlers.UserSessionHandlerFunc(c.GetCart))

		// route POST /api/v1/cart/items
		r.Method(http.MethodPost, "/items", handlers.UserSessionHandlerFunc(c.AddToCart))

		// route PATCH /api/v1/cart/items/{item_id}
		r.Method(http.MethodPatch, "/items/"+params.CartItemIdParam, handlers.UserSessionHandler(
			sfHandlers.CartItemHandler(c.UpdateCartItem)))

		// route DELETE /api/v1/cart/items/{item_id}
		r.Method(http.MethodDelete, "/items/"+params.CartItemIdParam, handlers.UserSessionHandler(
			sfHandlers.CartItemHandler(c.DeleteCartItem)))

	})
}
