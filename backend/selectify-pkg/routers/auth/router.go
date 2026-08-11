package auth

import (
	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/selectify-pkg/handlers"
)

func Register(r chi.Router) {
	c := new(controller).init()

	r.Route("/auth", func(r chi.Router) {
		// Post UserRegister /api/v1/auth/register
		r.Post("/register", c.UserRegister)
		// Post UserLogin /api/v1/auth/login
		r.Post("/login", c.UserLogin)

		// Post UserForgotPassword /api/v1/auth/forgot-password
		r.Post("/forgot-password", c.ForgotPassword)
		// Get ValidateResetPassword /api/v1/auth/reset-password/validate
		r.Get("/reset-password/validate", c.ValidateResetPassword)
		// Post UserResetPassword /api/v1/auth/reset-password
		r.Post("/reset-password", c.ResetPassword)
	})

	// Post UserLogout /api/v1/logout
	r.Post("/logout", handlers.UserSessionHandlerFunc(c.Logout))
}
