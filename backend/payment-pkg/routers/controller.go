package routers

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/payment-pkg/app"
)

type controller struct {
	paymentService service.PaymentService
}

func (c *controller) init() *controller {
	c.paymentService = app.Service().PaymentService

	return c
}
