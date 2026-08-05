package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type CheckoutRepo interface {
}

type checkoutRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewCheckoutRepo(db *db.DatabaseConnection) CheckoutRepo {
	return &checkoutRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (cr *checkoutRepo) CreateCheckoutSession(ctx context.Context, tx *sqlx.Tx, cs *model.CheckoutSession) error {
	q := `INSERT INTO checkout_session (order_id, status, expires_at) 
			VALUES (:order_id, :status, :expires_at) RETURNING id;`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		return logger.Error(ctx, err, "failed to prepare checkout session insert")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			_ = logger.Error(ctx, err, "failed to close statement")
		}
	}()

	if err = stmt.GetContext(ctx, &cs.ID, cs); err != nil {
		return logger.Error(ctx, err, "failed to insert checkout session")
	}

	return nil
}
