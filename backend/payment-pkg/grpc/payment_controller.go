package paymentgrpc

import (
	"context"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/service"
	"alwis.dev/selectify/payment-pkg/app"
	paymentv1 "alwis.dev/selectify/payment-pkg/grpc-proto/payment"
)

type controller struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func (c *controller) init() *controller {
	c.paymentService = app.Service().PaymentService

	return c
}

func NewController() paymentv1.PaymentServiceServer {
	return new(controller).init()
}

func (c *controller) CreatePayment(ctx context.Context, req *paymentv1.PaymentRequest) (*paymentv1.PaymentResponse, error) {
	if req.OrderId == 0 {
		return nil, logger.Error(ctx, nil, "orderId is required")
	}

	intent, err := c.paymentService.CreateStripeClient(ctx, req.OrderId)
	if err != nil {
		return nil, logger.Error(ctx, err, "error creating intent")
	}

	payment := &model.Payment{
		OrderID:           req.OrderId,
		Provider:          "stripe",
		ProviderPaymentID: intent.ID,
		ClientSecret:      intent.ClientSecret,
		Status:            string(intent.Status),
		Amount:            intent.Amount,
		Currency:          string(intent.Currency),
	}

	if err = c.paymentService.CreatePayment(ctx, payment); err != nil {
		return nil, logger.Error(ctx, err, "error creating payment")
	}

	return &paymentv1.PaymentResponse{
		Payment: paymentToProto(payment),
	}, nil
}

func paymentToProto(p *model.Payment) *paymentv1.Payment {
	if p == nil {
		return nil
	}

	out := &paymentv1.Payment{
		Id:                p.ID,
		OrderId:           p.OrderID,
		Provider:          p.Provider,
		ProviderPaymentId: p.ProviderPaymentID,
		ClientSecret:      p.ClientSecret,
		Status:            p.Status,
		Amount:            p.Amount,
		Currency:          p.Currency,
		CreatedAtUnix:     p.CreatedAt.Unix(),
		UpdatedAtUnix:     p.UpdatedAt.Unix(),
	}

	if p.FailureCode != nil {
		out.FailureCode = *p.FailureCode
	}
	if p.FailureMessage != nil {
		out.FailureMessage = *p.FailureMessage
	}
	if p.PaidAt != nil {
		out.PaidAtUnix = p.PaidAt.Unix()
	}
	if p.CreatedAt.IsZero() {
		out.CreatedAtUnix = 0
	}
	if p.UpdatedAt.IsZero() {
		out.UpdatedAtUnix = 0
	}

	return out
}
