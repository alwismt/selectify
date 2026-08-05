package service

import (
	"context"
	"log"
	"math"
	"os"

	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/paymentintent"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type PaymentService interface {
	CreateStripeClient(ctx context.Context, orderId uint64) (*stripe.PaymentIntent, error)
	CreatePayment(ctx context.Context, payment *model.Payment) error
	UpdatePaymentStatus(ctx context.Context, event stripe.Event) error
	//GetOrders(ctx context.Context, user *model.User) (model.Orders, error)
	//CreateOrder(ctx context.Context, user *model.User) (*model.Order, error)
}

type paymentService struct {
	stripeKey string
	//cartService CartService

	orderRepo   repo.OrderRepo
	paymentRepo repo.PaymentRepo
	//productVariantRepo repo.ProductVariantsRepo
	//txRepo             repo.TransactionRepo
}

func NewPaymentService(orderRepo repo.OrderRepo, paymentRepo repo.PaymentRepo) PaymentService {
	sKey := os.Getenv("PAY_STRIPE_SECRET_KEY")
	if sKey == "" {
		log.Fatal("payment_service.NewPaymentService: PAY_STRIPE_SECRET_KEY is not set")
	}

	return &paymentService{
		stripeKey:   sKey,
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
	}
}

func (o *paymentService) CreateStripeClient(ctx context.Context, orderId uint64) (*stripe.PaymentIntent, error) {
	order, err := o.orderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, logger.Error(ctx, err, "error getting order")
	}

	stripe.Key = o.stripeKey
	intent, err := paymentintent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(math.Round(order.Total * 100))),
		Currency: stripe.String(order.Currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	})

	if err != nil {
		return nil, logger.Error(ctx, err, "failed to get create Stripe client")
	}
	return intent, nil
}

func (o *paymentService) CreatePayment(ctx context.Context, payment *model.Payment) error {
	return o.paymentRepo.CreatePayment(ctx, payment)
}

func (o *paymentService) UpdatePaymentStatus(ctx context.Context, event stripe.Event) error {
	return logger.Error(ctx, nil, "not implemented")
}
