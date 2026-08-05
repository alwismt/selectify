package auth

import (
	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/handlers"
)

func Register(r chi.Router) {
	c := new(controller).init()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", c.UserRegister)
		r.Post("/login", c.UserLogin)
	})

	r.Post("/logout", handlers.UserSessionHandlerFunc(c.Logout))
}
