package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type CartRepo interface {
	UpsertCart(ctx context.Context, cart *model.Cart) error
	UpsertCartItem(ctx context.Context, cartItem *model.CartItem) error
	GetCartItemByID(ctx context.Context, itemID uint) (*model.CartItem, error)
	GetCartByID(ctx context.Context, userID, cartItemId uint) (*model.Cart, error)
	DeleteCartItem(ctx context.Context, itemID uint) error
	GetCartItemsByUserID(ctx context.Context, userID uint) (model.CartItems, error)
}

type cartRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewCartRepo(db *db.DatabaseConnection) CartRepo {
	return &cartRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *cartRepo) UpsertCart(ctx context.Context, cart *model.Cart) error {
	q := `INSERT INTO cart(user_id)
			VALUES ($1)
			ON CONFLICT (user_id)
			DO UPDATE SET updated_at = now()
			RETURNING id;`

	return r.rwDb.QueryRowContext(ctx, q, cart.UserID).Scan(&cart.ID)
}

func (r *cartRepo) UpsertCartItem(ctx context.Context, cartItem *model.CartItem) error {
	q := `INSERT INTO cart_items(cart_id, variant_id, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (cart_id, variant_id)
			DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = now();`

	_, err := r.rwDb.ExecContext(ctx, q, cartItem.CartID, cartItem.VariantID, cartItem.Quantity)
	if err != nil {
		err = logger.Errorf(ctx, err, "Failed to upsert cart item %d", cartItem.CartID)
		return err
	}

	return nil
}

func (r *cartRepo) GetCartItemByID(ctx context.Context, itemID uint) (*model.CartItem, error) {
	const q = `SELECT id, cart_id, variant_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE id = $1;`

	var cartItem model.CartItem
	if err := r.roDb.GetContext(ctx, &cartItem, q, itemID); err != nil {
		err = logger.Errorf(ctx, err, "Failed to get cart item by ID: %d", itemID)
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepo) GetCartByID(ctx context.Context, userID, cartID uint) (*model.Cart, error) {
	q := `SELECT id, user_id FROM cart WHERE id = $1 AND user_id = $2;`

	cart := new(model.Cart)
	if err := r.roDb.GetContext(ctx, cart, q, cartID, userID); err != nil {
		err = logger.Errorf(ctx, err, "Failed to get cart by ID: %d", cartID)
		return nil, err
	}
	return cart, nil
}

func (r *cartRepo) DeleteCartItem(ctx context.Context, itemID uint) error {
	q := `DELETE FROM cart_items WHERE id = $1;`
	res, err := r.rwDb.ExecContext(ctx, q, itemID)
	if err != nil {
		return logger.Errorf(ctx, err, "Failed to delete cart item %d", itemID)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return logger.Errorf(ctx, err, "Failed to delete cart item %d", itemID)
	}

	if rows == 0 {
		return logger.Error(ctx, sql.ErrNoRows, "")
	}

	return nil
}

func (r *cartRepo) GetCartItemsByUserID(ctx context.Context, userID uint) (model.CartItems, error) {
	q := `SELECT
			ci.id,
			ci.cart_id,
			ci.variant_id,
			ci.quantity,
			ci.created_at,
			ci.updated_at
		FROM cart_items ci
		JOIN cart c ON c.id = ci.cart_id
		WHERE c.user_id = $1
		ORDER BY ci.created_at;`

	var items model.CartItems
	if err := r.roDb.SelectContext(ctx, &items, q, userID); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get cart items for user %d", userID)
	}

	return items, nil
}
