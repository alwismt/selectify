package product

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	productService         service.ProductService
	productVariantsService service.ProductVariantsService
}

func (c *controller) init() *controller {
	c.productService = app.Service().ProductService
	c.productVariantsService = app.Service().ProductVariantsService

	return c
}
