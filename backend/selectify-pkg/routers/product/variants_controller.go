package product

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

func (c *controller) GetVariantsForProduct(w http.ResponseWriter, r *http.Request, product *model.Product) {
	ctx := r.Context()

	variants, err := c.productVariantsService.GetProVariantsForProduct(ctx, product)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to get product variants")
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err = httpx.NewJsonSender(variants, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(ctx, err, "failed to send response")
		httpx.SendError(w, httpx.ErrBadRequest)
	}
}
