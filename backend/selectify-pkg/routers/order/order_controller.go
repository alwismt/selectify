package order

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/httpx/request"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/params"
)

func (c *controller) CreateOrder(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	ctx := r.Context()

	order, err := c.orderService.CreateOrder(ctx, s.User)
	if err != nil {
		_ = logger.Error(ctx, err, "error creating order")
		httpx.SendError(w, fmt.Errorf("order creaing failed"))
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, order); err != nil {
		_ = logger.Error(ctx, err, "error sending response")
		httpx.SendError(w, err)
		return
	}

	return
}

func (c *controller) GetOrders(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	ctx := r.Context()

	orders, err := c.orderService.GetOrders(ctx, s.User)
	if err != nil {
		_ = logger.Error(ctx, err, "error getting orders")
		httpx.SendError(w, err)
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, orders); err != nil {
		_ = logger.Error(ctx, err, "error sending response")
		httpx.SendError(w, err)
		return
	}

	return
}

func (c *controller) SetOrderShippingAddress(w http.ResponseWriter, r *http.Request, s *model.LoggedInSession) {
	ctx := r.Context()

	orderIDStr := chi.URLParam(r, params.OrderId)
	if orderIDStr == "" {
		httpx.SendError(w, fmt.Errorf("order_id is required"))
		return
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil || orderID == 0 {
		httpx.SendError(w, fmt.Errorf("order_id is invalid"))
		return
	}

	req := new(request.OrderShippingAddressReq)
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

	addr := &model.OrderAddress{
		Phone:       req.Phone,
		Line1:       req.Line1,
		Line2:       req.Line2,
		City:        req.City,
		Region:      req.Region,
		PostalCode:  req.PostalCode,
		CountryCode: strings.ToUpper(req.CountryCode),
	}

	saved, err := c.orderService.SetOrderShippingAddress(ctx, s.User, uint(orderID), addr)
	if err != nil {
		_ = logger.Error(ctx, err, "error setting order shipping address")
		httpx.SendError(w, fmt.Errorf("failed to save shipping address"))
		return
	}

	if err = httpx.SendJson(w, http.StatusOK, saved); err != nil {
		_ = logger.Error(ctx, err, "error sending response")
		httpx.SendError(w, err)
		return
	}
}
