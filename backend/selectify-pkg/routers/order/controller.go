package order

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	orderService service.OrderService
}

func (c *controller) init() *controller {
	c.orderService = app.Service().OrderService

	return c
}
