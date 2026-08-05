package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"alwis.dev/selectify/internal/model"
	paymentv1 "alwis.dev/selectify/payment-pkg/grpc-proto/payment"
)

type PaymentClient struct {
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) PaymentClient {
	return PaymentClient{
		client: paymentv1.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentClient) CreatePayment(ctx context.Context, orderID uint) (*model.Payment, error) {
	res, err := c.client.CreatePayment(ctx, &paymentv1.PaymentRequest{OrderId: uint64(orderID)})
	if err != nil {
		return nil, err
	}

	return paymentFromProto(res.GetPayment()), nil
}

func paymentFromProto(p *paymentv1.Payment) *model.Payment {
	if p == nil {
		return nil
	}

	out := &model.Payment{
		ID:                p.Id,
		OrderID:           p.OrderId,
		Provider:          p.Provider,
		ProviderPaymentID: p.ProviderPaymentId,
		ClientSecret:      p.ClientSecret,
		Status:            p.Status,
		Amount:            p.Amount,
		Currency:          p.Currency,
	}

	if p.FailureCode != "" {
		code := p.FailureCode
		out.FailureCode = &code
	}
	if p.FailureMessage != "" {
		msg := p.FailureMessage
		out.FailureMessage = &msg
	}
	if p.CreatedAtUnix != 0 {
		out.CreatedAt = time.Unix(p.CreatedAtUnix, 0)
	}
	if p.UpdatedAtUnix != 0 {
		out.UpdatedAt = time.Unix(p.UpdatedAtUnix, 0)
	}
	if p.PaidAtUnix != 0 {
		t := time.Unix(p.PaidAtUnix, 0)
		out.PaidAt = &t
	}

	return out
}
