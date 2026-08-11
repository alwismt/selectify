package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/merchant-pkg/handlers"
)

//GET    /api/v1/merchant/dashboard
//
//GET    /api/v1/merchant/products
//POST   /api/v1/merchant/products
//GET    /api/v1/merchant/products/{id}
//PATCH  /api/v1/merchant/products/{id}
//DELETE /api/v1/merchant/products/{id}
//
//POST   /api/v1/merchant/products/{id}/images
//DELETE /api/v1/merchant/products/{id}/images/{imageId}

func RegisterRoutes(r chi.Router) {
	c := new(controller).init()

	r.Route("/products", func(r chi.Router) {
		r.Method(http.MethodGet, "/", handlers.MerchantSessionHandlerFunc(c.GetProduct))
	})
}
