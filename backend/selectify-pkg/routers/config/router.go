package config

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	//c := new(controller).init()
	r.Route("/config", func(r chi.Router) {
		// route GET /api/v1/config
		r.MethodFunc(http.MethodGet, "/", SiteConfig)
	})
}
