package merchant

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/merchant-pkg/app"
)

type controller struct {
	merchantService service.MerchantService
}

func (c *controller) init() *controller {
	c.merchantService = app.Service().MerchantService

	return c
}
