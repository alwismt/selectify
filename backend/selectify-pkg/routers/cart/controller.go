package cart

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	cartService service.CartService
}

func (c *controller) init() *controller {
	c.cartService = app.Service().CartService

	return c
}
