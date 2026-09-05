package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/params"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/selectify-pkg/app"
)

type ProductHandlerFunc func(http.ResponseWriter, *http.Request, *model.Product)

type productHandler struct {
	productRepo repo.ProductRepo

	handlerFunc ProductHandlerFunc
}

func (h productHandler) init(fn ProductHandlerFunc) productHandler {
	h.productRepo = app.Repository().ProductRepo
	h.handlerFunc = fn

	return h
}

func ProductHandler(fn ProductHandlerFunc) http.Handler {
	return productHandler{}.init(fn)
}

func (h productHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, params.ProductId)
	id, err := strconv.ParseUint(idStr, 10, strconv.IntSize)
	if err != nil {
		logger.Warnf(ctx, "failed to parse productID: %s", idStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.productRepo.GetProductByID(ctx, uint(id))
	if err != nil {
		logger.Warnf(ctx, "failed to get product by ID: %d: %s", id, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.handlerFunc(w, r, product)
}
