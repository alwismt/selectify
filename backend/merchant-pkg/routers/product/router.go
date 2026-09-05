package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/params"
	"alwis.dev/selectify/internal/permission"
	"alwis.dev/selectify/merchant-pkg/handlers"
)

//GET    /api/v1/merchant/dashboard
//
//GET    /api/v1/merchant/products
//POST   /api/v1/merchant/products
//GET    /api/v1/merchant/products/{id}
//PATCH  /api/v1/merchant/products/{id}
//DELETE /api/v1/merchant/products/{id}
//
//POST   /api/v1/merchant/products/{id}/images
//DELETE /api/v1/merchant/products/{id}/images/{imageId}

func RegisterRoutes(r chi.Router) {
	c := new(controller).init()
	r.Route("/products", func(r chi.Router) {
		// GET /api/v1/merchant/products GetProducts
		r.Method(http.MethodGet, "/",
			handlers.MerchantSessionHandler(
				handlers.MerchantAccessHandler(c.GetProducts, permission.Product.Read)))
		// GET /api/v1/merchant/products/{product_id} GetProductByID
		r.Method(http.MethodGet, "/"+params.ProductIdParm,
			handlers.MerchantSessionHandler(
				handlers.MerchantAccessHandler(
					handlers.MerchantProductHandler(c.GetProductByID),
					permission.Product.Read)))
		// GET /api/v1/merchant/products/{product_id}/variants GetVariantsForProduct
		r.Method(http.MethodGet, "/"+params.ProductIdParm+"/variants",
			handlers.MerchantSessionHandler(
				handlers.MerchantAccessHandler(
					handlers.MerchantProductHandler(c.GetVariantsForProduct),
					permission.Product.Read)))

		// POST /api/v1/merchant/products CreateProduct
		r.Method(http.MethodPost, "/", handlers.MerchantSessionHandler(
			handlers.MerchantAccessHandler(c.CreateProduct, permission.Product.Create)))
	})
}
