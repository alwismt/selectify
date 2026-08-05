package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/params"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/selectify-pkg/app"
)

type CartItemHandlerFunc func(http.ResponseWriter, *http.Request, *model.UserSession, *model.CartItem)

type cartItemHandler struct {
	cartRepo repo.CartRepo

	handlerFunc CartItemHandlerFunc
}

func (h cartItemHandler) init(fn CartItemHandlerFunc) cartItemHandler {
	h.cartRepo = app.Repository().CartRepo

	h.handlerFunc = fn

	return h
}

func CartItemHandler(fn CartItemHandlerFunc) cartItemHandler {
	return cartItemHandler{}.init(fn)
}

func (h cartItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	ctx := r.Context()
	cartItemIdStr := chi.URLParam(r, params.CartItemId)
	if cartItemIdStr == "" {
		logger.Warn(ctx, "cartItemId is empty")
		httpx.SendError(w, fmt.Errorf("cartItemId is required"))
		return
	}

	cartItemId, err := strconv.Atoi(cartItemIdStr)
	if err != nil {
		logger.Warn(ctx, "cartItemId is invalid")
		httpx.SendError(w, fmt.Errorf("cartItemId is invalid"))
		return
	}

	cart, err := h.cartRepo.GetCartItemByID(ctx, uint(cartItemId))
	if err != nil {
		logger.Warn(ctx, "cartItemId is invalid")
		httpx.SendError(w, fmt.Errorf("cartItemId is invalid"))
		return
	}

	if cart == nil {
		logger.Warn(ctx, "cartItemId is nil")
		httpx.SendError(w, fmt.Errorf("cartItemId is nil"))
		return
	}

	cart.Cart, err = h.cartRepo.GetCartByID(ctx, s.UserId, cart.CartID)
	if err != nil {
		logger.Warn(ctx, "cartItemId is invalid")
		httpx.SendError(w, fmt.Errorf("cartItemId is invalid"))
		return
	}

	h.handlerFunc(w, r, s, cart)
}
