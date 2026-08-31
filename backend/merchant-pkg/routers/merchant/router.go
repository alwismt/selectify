package merchant

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/permission"
	"alwis.dev/selectify/merchant-pkg/handlers"
)

// GET    /api/v1/merchant
// Get the logged-in user's merchant/store
// POST   /api/v1/merchant
// Become a merchant / create store
// PATCH  /api/v1/merchant
// Update store details

func RegisterRoutes(r chi.Router) {
	c := new(controller).init()
	// GET /api/v1/merchant GetMerchant
	r.Method(http.MethodGet, "/", handlers.MerchantSessionHandler(handlers.MerchantAccessHandler(
		c.GetMerchant, permission.Merchant.Read)))
	// PATCH /api/v1/merchant UpdateMerchant
	r.Method(http.MethodPatch, "/", handlers.MerchantSessionHandler(handlers.MerchantAccessHandler(
		c.UpdateMerchant, permission.Merchant.Update)))
	// PATCH /api/v1/merchant/logo UpdateLogoMerchant
	r.Method(http.MethodPost, "/logo", handlers.MerchantSessionHandler(handlers.MerchantAccessHandler(
		c.UpdateMerchantLogo, permission.Merchant.Update)))

	// GET /api/v1/merchant/countries GetMerchant Countries
	r.MethodFunc(http.MethodGet, "/countries", handlers.MerchantSessionHandlerFunc(c.GetMerchantCountries))
}
