package status

import (
	"github.com/go-chi/chi/v5"
)

func Register(r chi.Router) {
	r.Get("/status", GetStatus)
}
