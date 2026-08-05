package user

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/selectify-pkg/app"
)

func UserInfo(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	if err := httpx.SendJson(w, http.StatusOK, s.User); err != nil {
		_ = logger.Error(r.Context(), err, "failed to send user info response")
	}

	return
}

func GetDefaultAddress(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	ctx := r.Context()

	addr, err := app.Repository().UserAddressRepo.GetDefaultByUserID(ctx, s.User.ID)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to get default user address")
		httpx.SendError(w, err)
		return
	}

	if addr == nil {
		if err = httpx.SendJson(w, http.StatusOK, map[string]any{}); err != nil {
			_ = logger.Error(ctx, err, "failed to send empty address response")
		}
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, addr); err != nil {
		_ = logger.Error(ctx, err, "failed to send user address response")
	}
}
