package service

import (
	"context"

	"alwis.dev/selectify/internal/httpx/response"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type CartService interface {
	AddToCart(ctx context.Context, variantId, qnt uint, user *model.User) error
	ItemUpsert(ctx context.Context, qnt uint, cartItem *model.CartItem) error
	DeleteCartItem(ctx context.Context, cartItem *model.CartItem) error
	GetCartItems(ctx context.Context, s *model.User) (*response.CartResponse, error)
}

type cartService struct {
	cartRepo     repo.CartRepo
	productRepo  repo.ProductRepo
	variantsRepo repo.ProductVariantsRepo
}

func NewCartService(cartRepo repo.CartRepo, productRepo repo.ProductRepo, variantsRepo repo.ProductVariantsRepo) CartService {
	return &cartService{
		cartRepo:     cartRepo,
		productRepo:  productRepo,
		variantsRepo: variantsRepo,
	}
}

func (cs *cartService) AddToCart(ctx context.Context, variantId, qnt uint, user *model.User) error {
	v, err := cs.variantsRepo.GetVariantsByID(ctx, variantId)
	if err != nil {
		err = logger.Errorf(ctx, err, "failed to get variant %d", variantId)
		return err
	}

	if v == nil {
		err = logger.Errorf(ctx, err, "failed to get variant %d", variantId)
		return err
	}

	if v.AvailableStockQty() < qnt {
		return logger.Errorf(ctx, err, "insufficient stock for variant %d (available=%d, requested=%d)", variantId, v.AvailableStockQty(), qnt)
	}

	cart := &model.Cart{
		UserID: user.ID,
	}

	if err = cs.cartRepo.UpsertCart(ctx, cart); err != nil {
		err = logger.Errorf(ctx, err, "failed to upsert cart")
		return err
	}

	attr := &model.CartItem{
		CartID:    cart.ID,
		VariantID: v.ID,
		Quantity:  qnt,
	}

	if err = cs.cartRepo.UpsertCartItem(ctx, attr); err != nil {
		err = logger.Errorf(ctx, err, "failed to upsert cart %d item for variant %d", attr.CartID, v.ID)
		return err
	}

	return nil
}

func (cs *cartService) ItemUpsert(ctx context.Context, qnt uint, cartItem *model.CartItem) error {
	v, err := cs.variantsRepo.GetVariantsByID(ctx, cartItem.VariantID)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to get variant %d", cartItem.VariantID)
	}

	if v.AvailableStockQty() < qnt {
		return logger.Errorf(ctx, err, "insufficient stock for variant %d (available=%d, requested=%d)", cartItem.VariantID, v.AvailableStockQty(), qnt)
	}

	cartItem.Quantity = qnt
	if err = cs.cartRepo.UpsertCartItem(ctx, cartItem); err != nil {
		err = logger.Errorf(ctx, err, "failed to upsert cart")
		return err
	}

	return nil
}

func (cs *cartService) DeleteCartItem(ctx context.Context, cartItem *model.CartItem) error {
	if err := cs.cartRepo.DeleteCartItem(ctx, cartItem.ID); err != nil {
		return logger.Errorf(ctx, err, "failed to delete cart item")
	}
	return nil
}

func (cs *cartService) GetCartItems(ctx context.Context, s *model.User) (*response.CartResponse, error) {
	crItems, err := cs.cartRepo.GetCartItemsByUserID(ctx, s.ID)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get cart items")
	}

	if len(crItems) == 0 {
		return &response.CartResponse{
			Items:     []response.CartItem{},
			Subtotal:  0,
			ItemCount: 0,
		}, nil
	}

	vrIds := crItems.GetVariantIDs()
	if len(vrIds) == 0 {
		return nil, logger.Errorf(ctx, err, "failed to get cart variants items")
	}

	variants, err := cs.variantsRepo.GetVariantsByIDs(ctx, vrIds)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get protocol variants for ids %v", vrIds)
	}

	ids := make([]uint, 0, len(variants))
	variantMap := map[uint]*model.ProductVariant{}

	for i := range variants {
		v := &variants[i]
		ids = append(ids, v.ID)
		variantMap[v.ID] = v
	}

	attrs, err := cs.variantsRepo.GetAttrByVariantsIDs(ctx, ids)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get protocol variants for ids %v", ids)
	}

	for _, a := range attrs {
		if v := variantMap[a.VariantID]; v != nil {
			v.ProductVariantAttributes = append(v.ProductVariantAttributes, a)
		}
	}

	productIDSet := make(map[uint]bool)
	productIDs := make([]uint, 0)
	for _, v := range variants {
		if !productIDSet[v.ProductID] {
			productIDSet[v.ProductID] = true
			productIDs = append(productIDs, v.ProductID)
		}
	}

	products, err := cs.productRepo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get products for ids %v", productIDs)
	}

	productMap := make(map[uint]*model.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}
	return cs.buildCartResponse(crItems, variantMap, productMap), nil
}

func (cs *cartService) buildCartResponse(crItems model.CartItems, variantMap map[uint]*model.ProductVariant, productMap map[uint]*model.Product,
) *response.CartResponse {
	items := make([]response.CartItem, 0, len(crItems))
	var subtotal uint64
	itemCount := uint(len(crItems))

	for _, cartItem := range crItems {
		variant := variantMap[cartItem.VariantID]
		if variant == nil {
			continue
		}

		product := productMap[variant.ProductID]
		if product == nil {
			continue
		}

		attributes := make(map[string]string)
		for _, attr := range variant.ProductVariantAttributes {
			attributes[attr.Name] = attr.Value
		}

		var itemPrice uint64
		if variant.PriceAmount != nil {
			itemPrice = *variant.PriceAmount
		}
		subtotal += itemPrice * uint64(cartItem.Quantity)

		availableQty := variant.StockQty - variant.ReservedQty
		respVariant := response.ProductVariant{
			ID:          variant.ID,
			SKU:         variant.SKU,
			PriceAmount: variant.PriceAmount,
			Attributes:  attributes,
			Product: response.Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
			},
			StockQty: availableQty,
		}

		items = append(items, response.CartItem{
			ID:       cartItem.ID,
			Quantity: cartItem.Quantity,
			Variant:  respVariant,
		})
	}

	return &response.CartResponse{
		Items:     items,
		Subtotal:  subtotal,
		ItemCount: itemCount,
	}
}
