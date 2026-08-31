package merchant

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/merchant-pkg/app"
)

type controller struct {
	countryService  service.CountryService
	merchantService service.MerchantService
}

func (c *controller) init() *controller {
	c.countryService = app.Service().CountryService
	c.merchantService = app.Service().MerchantService

	return c
}
