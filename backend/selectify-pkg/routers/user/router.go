package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	//"alwis.dev/selectify/selectify-pkg/handlers"

	//"alwis.dev/selectify/selectify-pkg/handlers"

	"alwis.dev/selectify/internal/handlers"
)

func Register(r chi.Router) {
	c := new(controller).init()

	r.Group(func(r chi.Router) {
		//sessionRepo := app.Repository().UserSessionRepo
		//r.Use(middleware.ContextMiddleware(sessionRepo))

		r.Route("/user", func(r chi.Router) {
			r.MethodFunc(http.MethodGet, "/info", handlers.UserSessionHandlerFunc(UserInfo))

			r.Get("/me", handlers.UserSessionHandlerFunc(c.GetUserImage))
			r.Post("/me", handlers.UserSessionHandlerFunc(c.UpdateUserImage))
			r.Delete("/me", handlers.UserSessionHandlerFunc(c.DeleteUserImage))

			r.MethodFunc(http.MethodGet, "/addresses/default", handlers.UserSessionHandlerFunc(GetDefaultAddress))
		})
	})
}
