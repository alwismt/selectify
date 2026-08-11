package merchant

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/merchant-pkg/handlers"
)

//GET    /api/v1/merchant
//Get the logged-in user's merchant/store
//POST   /api/v1/merchant
//Become a merchant / create store
//PATCH  /api/v1/merchant
//Update store details

func RegisterRoutes(r chi.Router) {
	c := new(controller).init()
	// GET    /api/v1/merchant
	r.Method(http.MethodGet, "/", handlers.MerchantSessionHandlerFunc(c.GetMerchant))
}
