package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type PaymentRepo interface {
	CreatePayment(ctx context.Context, payment *model.Payment) error
}

type paymentRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewPaymentRepo(db *db.DatabaseConnection) PaymentRepo {
	return &paymentRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *paymentRepo) CreatePayment(ctx context.Context, payment *model.Payment) error {
	q := `INSERT INTO payments(order_id, provider, provider_payment_id, client_secret, status, amount, currency) 
		VALUES (:order_id, :provider, :provider_payment_id, :client_secret, :status, :amount, :currency)
		RETURNING id;`

	rows, err := r.rwDb.NamedQueryContext(ctx, q, payment)
	if err != nil {
		return logger.Error(ctx, err, "error creating payment")
	}
	defer rows.Close()

	if !rows.Next() {
		return logger.Error(ctx, sql.ErrNoRows, "payment insert returned no ID")
	}
	if err = rows.Scan(&payment.ID); err != nil {
		return logger.Error(ctx, err, "error scanning payment ID")
	}

	return nil
}

func (r *paymentRepo) UpdatePaymentStatus(ctx context.Context, payment *model.Payment) error {
	return logger.Error(ctx, nil, "not implemented")
}
