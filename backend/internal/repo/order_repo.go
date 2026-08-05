package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type OrderRepo interface {
	GetOrderById(ctx context.Context, orderId uint64) (*model.Order, error)
	GetOrders(ctx context.Context, userID uint) (*model.Orders, error)
	LoadOrderItems(ctx context.Context, orders *model.Orders) error
	CreateOrder(ctx context.Context, tx *sqlx.Tx, order *model.Order) error
	CreateOrderItems(ctx context.Context, tx *sqlx.Tx, items model.OrderItems) error
	UpsertOrderAddress(ctx context.Context, addr *model.OrderAddress) error
}

type orderRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewOrderRepo(db *db.DatabaseConnection) OrderRepo {
	return &orderRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (or *orderRepo) CreateOrder(ctx context.Context, tx *sqlx.Tx, order *model.Order) error {
	q := `INSERT INTO orders (user_id, status, currency, subtotal, shipping, discount, total)
		VALUES ( :user_id, :status, :currency, :subtotal, :shipping, :discount, :total)
		RETURNING id;`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		return logger.Error(ctx, err, "failed to prepare order insert")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			_ = logger.Error(ctx, err, "failed to close statement")
		}
	}()

	if err = stmt.GetContext(ctx, &order.ID, order); err != nil {
		return logger.Error(ctx, err, "failed to insert order")
	}

	return nil
}

func (or *orderRepo) CreateOrderItems(ctx context.Context, tx *sqlx.Tx, items model.OrderItems) error {
	if len(items) == 0 {
		return nil
	}

	base := `INSERT INTO order_items (
		order_id, variant_id, sku, unit_price, currency, quantity, attributes
	) VALUES %s RETURNING id, created_at;`

	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]any, 0, len(items)*7)

	for i, it := range items {
		n := i * 7
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				n+1, n+2, n+3, n+4, n+5, n+6, n+7,
			),
		)

		valueArgs = append(valueArgs,
			it.OrderID,
			it.VariantID,
			it.SKU,
			it.UnitPrice,
			it.Currency,
			it.Quantity,
			it.Attributes,
		)
	}

	q := fmt.Sprintf(base, strings.Join(valueStrings, ","))

	rows, err := tx.QueryxContext(ctx, q, valueArgs...)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to batch insert order items (count=%d)", len(items))
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		if err := rows.Scan(&items[i].ID, &items[i].CreatedAt); err != nil {
			return logger.Error(ctx, err, "failed to scan returned order item")
		}
		i++
	}

	return nil
}

func (or *orderRepo) GetOrders(ctx context.Context, userID uint) (*model.Orders, error) {
	var orders model.Orders

	const q = `SELECT * FROM orders WHERE user_id = $1 ORDER BY id DESC;`
	if err := or.rwDb.SelectContext(ctx, &orders, q, userID); err != nil {
		return nil, logger.Error(ctx, err, "failed to get orders")
	}

	return &orders, nil
}

func (or *orderRepo) LoadOrderItems(ctx context.Context, orders *model.Orders) error {
	if orders == nil || len(*orders) == 0 {
		return nil
	}

	ids := make([]uint, 0, len(*orders))
	orderIndex := make(map[uint]int, len(*orders))

	for i := range *orders {
		id := (*orders)[i].ID
		ids = append(ids, id)
		orderIndex[id] = i
		(*orders)[i].Items = nil
	}

	q, args, err := sqlx.In(`
		SELECT *
		FROM order_items
		WHERE order_id IN (?)
		ORDER BY order_id DESC, id ASC;`, ids)

	if err != nil {
		return logger.Error(ctx, err, "failed to build order items query")
	}

	q = or.rwDb.Rebind(q)

	var items model.OrderItems
	if err := or.rwDb.SelectContext(ctx, &items, q, args...); err != nil {
		return logger.Error(ctx, err, "failed to get order items")
	}

	for _, it := range items {
		if idx, ok := orderIndex[it.OrderID]; ok {
			(*orders)[idx].Items = append((*orders)[idx].Items, it)
		}
	}

	return nil
}

func (or *orderRepo) GetOrderById(ctx context.Context, orderId uint64) (*model.Order, error) {
	var order model.Order
	if err := or.rwDb.GetContext(ctx, &order, `SELECT * FROM orders WHERE id = $1`, orderId); err != nil {
		return nil, logger.Error(ctx, err, "failed to get order")
	}
	return &order, nil
}

func (or *orderRepo) UpsertOrderAddress(ctx context.Context, addr *model.OrderAddress) error {
	const q = `
		INSERT INTO order_addresses (
			order_id, type, phone, line1, line2, city, region, postal_code, country_code
		) VALUES (
			:order_id, :type, :phone, :line1, :line2, :city, :region, :postal_code, :country_code
		)
		ON CONFLICT (order_id, type) DO UPDATE SET
			phone = EXCLUDED.phone,
			line1 = EXCLUDED.line1,
			line2 = EXCLUDED.line2,
			city = EXCLUDED.city,
			region = EXCLUDED.region,
			postal_code = EXCLUDED.postal_code,
			country_code = EXCLUDED.country_code
		RETURNING id, created_at;`

	stmt, err := or.rwDb.PrepareNamedContext(ctx, q)
	if err != nil {
		return logger.Error(ctx, err, "failed to prepare order address upsert")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			_ = logger.Error(ctx, closeErr, "failed to close statement")
		}
	}()

	if err = stmt.GetContext(ctx, addr, addr); err != nil {
		return logger.Error(ctx, err, "failed to upsert order address")
	}
	return nil
}
