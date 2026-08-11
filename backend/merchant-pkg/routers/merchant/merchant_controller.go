package merchant

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

func (c *controller) GetMerchant(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	merchant, err := c.merchantService.GetMerchant(r.Context(), *s.UserRole.MerchantID)
	if err != nil {
		err = logger.Error(r.Context(), err, "Failed to get merchant")
		httpx.SendUnAuthorized(w)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, merchant); err != nil {
		_ = logger.Error(r.Context(), err, "Failed to send response")
	}
	return
}
