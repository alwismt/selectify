package config

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/selectify-pkg/app"
)

func SiteConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := app.Service().SiteConfigService

	config, err := s.GetSiteConfig(ctx)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to get site config")
		httpx.SendError(w, httpx.ErrInternalServer)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, config); err != nil {
		_ = logger.Error(ctx, err, "Failed to send site config")
		httpx.SendError(w, httpx.ErrInternalServer)
	}
	return
}
