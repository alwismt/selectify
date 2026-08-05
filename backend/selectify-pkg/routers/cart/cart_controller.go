package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

func (c *controller) AddToCart(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	ctx := r.Context()
	req := new(request.AddCartReq)

	if err := httpx.MustDecodeJson(w, r, req); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to decode body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		fieldErrors := httpx.ValidationErrorsToMap(err)
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	if err := c.cartService.AddToCart(ctx, req.VariantId, req.Quantity, s.User); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to add to cart")
		_ = httpx.SendJson(w, http.StatusBadRequest, "Failed to add to cart")
		return
	}

	if err := httpx.SendJson(w, http.StatusCreated, "success"); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to send success")
	}
	return
}

func (c *controller) UpdateCartItem(w http.ResponseWriter, r *http.Request, _ *model.UserSession, cart *model.CartItem) {
	ctx := r.Context()

	req := new(request.QntReq)
	if err := httpx.MustDecodeJson(w, r, req); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to decode body")
		return
	}

	if err := validator.New().Struct(req); err != nil {
		_ = logger.Errorf(ctx, err, "Validation error")
		fieldErrors := httpx.ValidationErrorsToMap(err)
		_ = httpx.SendJson(w, http.StatusBadRequest, fieldErrors)
		return
	}

	if req.Quantity == cart.Quantity {
		httpx.SendNoContent(w)
		return
	}

	if err := c.cartService.ItemUpsert(ctx, req.Quantity, cart); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to update item")
		httpx.SendError(w, fmt.Errorf("failed to update item"))
		return
	}

	if err := httpx.SendJson(w, http.StatusOK, "success"); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to send success")
	}
	return
}

func (c *controller) DeleteCartItem(w http.ResponseWriter, r *http.Request, _ *model.UserSession, cart *model.CartItem) {
	ctx := r.Context()

	if err := c.cartService.DeleteCartItem(ctx, cart); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to delete item")
		httpx.SendError(w, fmt.Errorf("failed to delete item"))
		return
	}
	httpx.SendNoContent(w)
	return
}

func (c *controller) GetCart(w http.ResponseWriter, r *http.Request, s *model.UserSession) {
	ctx := r.Context()

	cartResponse, err := c.cartService.GetCartItems(ctx, s.User)
	if err != nil {
		_ = logger.Errorf(ctx, err, "Failed to get items")
		httpx.SendError(w, fmt.Errorf("failed to get items"))
		return
	}

	if err := httpx.SendJson(w, http.StatusOK, cartResponse); err != nil {
		_ = logger.Errorf(ctx, err, "Failed to send cart response")
		return
	}

	return
}
