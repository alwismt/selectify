package product

import (
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/merchant-pkg/app"
)

type controller struct {
	productService         service.ProductService
	productFileService     service.ProductFileService
	productVariantsService service.ProductVariantsService
}

func (c *controller) init() *controller {
	c.productService = app.Service().ProductService
	c.productFileService = app.Service().ProductFileService
	c.productVariantsService = app.Service().ProductVariantsService
	return c
}
