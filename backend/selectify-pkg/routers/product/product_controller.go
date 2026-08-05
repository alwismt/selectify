package product

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/params"
)

func (c *controller) GetProductById(w http.ResponseWriter, r *http.Request, product *model.Product) {
	ctx := r.Context()
	var err error
	product, err = c.productService.GetProductFileByProductID(ctx, product)
	if err != nil {
		_ = logger.Error(ctx, err, "failed to fetch product file")
		httpx.SendError(w, err)
		return
	}

	if err := httpx.NewJsonSender(product, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(ctx, err, "failed to send response")
		httpx.SendError(w, err)
	}
	return
}

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

func (c *controller) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	slug := chi.URLParam(r, params.ProductSlug)
	if slug == "" {
		err := logger.Error(ctx, errors.New("failed to get product url"), "url is required")
		httpx.SendError(w, err)
		return
	}

	// validate slug
	if !slugRegex.MatchString(slug) {
		err := logger.Error(ctx, errors.New("invalid url format"), "invalid url")
		httpx.SendError(w, err)
		return
	}

	product, err := c.productService.GetProductBySlug(ctx, slug)
	if err != nil {
		err = logger.Error(ctx, err, "failed to fetch product")
		httpx.SendError(w, err)
		return
	}

	if err = httpx.NewJsonSender(product, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(ctx, err, "failed to send response")
		httpx.SendError(w, err)
	}
}

// GetProducts Todo:: ?search=&category_id=&min_price=&max_price=&sort=&page=&limit=
func (c *controller) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := c.productService.GetProducts(ctx)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to get products")
		httpx.SendError(w, err)
		return
	}

	if err = httpx.NewJsonSender(products, http.StatusOK).Send(w); err != nil {
		_ = logger.Error(ctx, err, "failed to send response")
		httpx.SendError(w, err)
	}

	return
}
