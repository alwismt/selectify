package category

import (
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	c := new(controller).init()

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", c.GetCategories)
	})
}
