package service

import (
	"context"
	"encoding/json"
	"fmt"

	grpcclient "alwis.dev/selectify/internal/grpc-client"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
	"alwis.dev/selectify/internal/types"
)

type OrderService interface {
	GetOrders(ctx context.Context, user *model.User) (model.Orders, error)
	CreateOrder(ctx context.Context, user *model.User) (*model.Order, error)
	SetOrderShippingAddress(ctx context.Context, user *model.User, orderID uint, addr *model.OrderAddress) (*model.OrderAddress, error)
}

type paymentClient interface {
	CreatePayment(ctx context.Context, orderID uint) (*model.Payment, error)
}

type orderService struct {
	cartService CartService

	paymentClient grpcclient.PaymentClient

	orderRepo          repo.OrderRepo
	userAddressRepo    repo.UserAddressRepo
	productVariantRepo repo.ProductVariantsRepo
	txRepo             repo.TransactionRepo
}

func NewOrderService(cartService CartService, orderRepo repo.OrderRepo, txRepo repo.TransactionRepo,
	variantsRepo repo.ProductVariantsRepo, paymentClient grpcclient.PaymentClient,
	userAddressRepo repo.UserAddressRepo) OrderService {
	return &orderService{
		cartService:        cartService,
		paymentClient:      paymentClient,
		orderRepo:          orderRepo,
		userAddressRepo:    userAddressRepo,
		productVariantRepo: variantsRepo,
		txRepo:             txRepo,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, user *model.User) (*model.Order, error) {
	cart, err := o.cartService.GetCartItems(ctx, user)
	if err != nil || cart == nil {
		return nil, logger.Error(ctx, err, "failed to get cart items")
	}

	order := &model.Order{
		UserID:   user.ID,
		Status:   types.OrderStatusNew,
		Currency: cart.Currency,
		Subtotal: cart.Subtotal,
		Shipping: 0,
		Discount: 0,
		Total:    cart.Subtotal,
	}

	// Todo:: setup checkout_session allowing user to checkout selected items from cart

	tx, err := o.txRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to start transaction")
	}
	defer tx.End()

	if err = o.orderRepo.CreateOrder(ctx, tx.Transaction, order); err != nil {
		return nil, logger.Error(ctx, err, "failed to create order")
	}

	//checkoutSession := &model.CheckoutSession{
	//	OrderID:   order.ID,
	//	Status:    types.PaymentStatusPending,
	//	ExpiresAt: time.Now().Add(time.Minute * 30),
	//}
	var (
		orderItems model.OrderItems
		//checkoutSessionItem model.CheckoutSessionItems
	)

	for _, i := range cart.Items {
		attrsJSON, err := json.Marshal(i.Variant.Attributes)
		if err != nil {
			return nil, logger.Error(ctx, err, "failed to marshal variant attributes")
		}

		// Todo:: Make this a bulk update
		if err = o.productVariantRepo.ReserveStockForCheckout(ctx, tx.Transaction, i.Variant.ID, i.Quantity); err != nil {
			return nil, logger.Error(ctx, err, "failed to reserve stock for checkout")
		}

		items := model.OrderItem{
			OrderID:    order.ID,
			VariantID:  i.Variant.ID,
			SKU:        i.Variant.SKU,
			UnitPrice:  *i.Variant.PriceAmount,
			Currency:   i.Variant.Currency,
			Quantity:   i.Quantity,
			Attributes: json.RawMessage(attrsJSON),
		}

		orderItems = append(orderItems, items)

		//sessionItem := model.CheckoutSessionItem{
		//	CheckoutSessionID: checkoutSession.ID,
		//	CartItemID:        i.ID,
		//}
		//
		//checkoutSessionItem = append(checkoutSessionItem, sessionItem)
	}

	if err = o.orderRepo.CreateOrderItems(ctx, tx.Transaction, orderItems); err != nil {
		return nil, logger.Error(ctx, err, "failed to create order items")
	}

	order.Items = orderItems
	tx.CanCommit = true
	tx.End()

	payment, err := o.paymentClient.CreatePayment(ctx, order.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to create payment")
	}
	order.ClientSecret = payment.ClientSecret

	return order, nil
}

func (o *orderService) GetOrders(ctx context.Context, user *model.User) (model.Orders, error) {
	orders, err := o.orderRepo.GetOrders(ctx, user.ID)
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to get orders")
	}

	if err = o.orderRepo.LoadOrderItems(ctx, orders); err != nil {
		return nil, logger.Error(ctx, err, "failed to load order items")
	}

	return *orders, nil
}

func (o *orderService) SetOrderShippingAddress(ctx context.Context, user *model.User, orderID uint, addr *model.OrderAddress) (*model.OrderAddress, error) {
	order, err := o.orderRepo.GetOrderById(ctx, uint64(orderID))
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to get order")
	}

	if order.UserID != user.ID {
		return nil, logger.Error(ctx, fmt.Errorf("order not found"), "order does not belong to user")
	}

	switch order.Status {
	case types.OrderStatusNew, types.OrderStatusPending:
		// payable
	default:
		return nil, logger.Errorf(ctx, fmt.Errorf("order is not payable"), "cannot set address for order status %s", order.Status)
	}

	addr.OrderID = order.ID
	addr.Type = "shipping"

	if err = o.orderRepo.UpsertOrderAddress(ctx, addr); err != nil {
		return nil, logger.Error(ctx, err, "failed to upsert order address")
	}

	userAddr := &model.UserAddress{
		UserID:      user.ID,
		Phone:       addr.Phone,
		Line1:       addr.Line1,
		Line2:       addr.Line2,
		City:        addr.City,
		Region:      addr.Region,
		PostalCode:  addr.PostalCode,
		CountryCode: addr.CountryCode,
		IsDefault:   true,
	}
	if err = o.userAddressRepo.UpsertDefaultAddress(ctx, userAddr); err != nil {
		return nil, logger.Error(ctx, err, "failed to upsert user address")
	}

	return addr, nil
}
