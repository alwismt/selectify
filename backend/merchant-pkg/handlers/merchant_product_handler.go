package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/params"
	"alwis.dev/selectify/merchant-pkg/app"
)

type MerchantProductHandlerFunc func(http.ResponseWriter, *http.Request, *model.LoggedInSession, *model.Merchant, *model.Product)

func MerchantProductHandler(fn MerchantProductHandlerFunc) MerchantAccessHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession, m *model.Merchant) {
		ctx := r.Context()

		productId := chi.URLParam(r, params.ProductId)
		if productId == "" {
			_ = logger.Error(ctx, errors.New("failed to get product url"), "url is required")
			httpx.SendError(w, httpx.ErrInvalidRequest)
			return
		}

		//idPart, _, found := strings.Cut(productPath, "-")
		//
		//if !found || idPart == "" {
		//	w.WriteHeader(http.StatusNotFound)
		//	return
		//}

		productID, err := strconv.ParseUint(productId, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		product, err := app.Repository().ProductRepo.GetMerchantProductByID(ctx, m.MerchantID, uint(productID))
		if err != nil {
			logger.Warnf(ctx, "failed to get product by slug %d, merchant ID: %d, %s", productID, m.MerchantID, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fn(w, r, s, m, product)
	}
}
