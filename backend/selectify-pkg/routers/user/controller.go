package user

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	userService service.UserService
}

func (c *controller) init() *controller {
	c.userService = app.Service().UserService
	return c
}
