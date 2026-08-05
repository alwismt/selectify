package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/handlers"
	"alwis.dev/selectify/internal/params"
)

func Register(r chi.Router) {
	c := new(controller).init()

	r.Route("/products", func(r chi.Router) {
		r.Get("/", c.GetProducts)
		r.Method(http.MethodGet, "/id/"+params.ProductIdParm, handlers.ProductHandler(c.GetProductById))
		r.MethodFunc(http.MethodGet, "/"+params.ProductSlugParm, c.GetProductBySlug)

		// variants
		r.Method(http.MethodGet, "/"+params.ProductIdParm+"/variants", handlers.ProductHandler(c.GetVariantsForProduct))
	})
}
