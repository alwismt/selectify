package category

import (
	"net/http"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/selectify-pkg/app"
)

type controller struct {
	categoryService service.CategoryService
}

func (c *controller) init() *controller {
	c.categoryService = app.Service().CategoryService
	return c
}

func (c *controller) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := c.categoryService.GetCategories(ctx)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to get categories")
		httpx.SendError(w, httpx.ErrBadRequest)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, categories); err != nil {
		_ = logger.Error(ctx, err, "failed to send categories")
		httpx.SendError(w, httpx.ErrBadRequest)
	}
	return
}
