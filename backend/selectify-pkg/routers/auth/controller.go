package auth

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	authService service.AuthService
}

func (c *controller) init() *controller {
	c.authService = app.Service().AuthService

	return c
}
