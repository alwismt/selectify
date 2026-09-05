package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type CurrencyRepo interface {
	GetDefaultCurrency(ctx context.Context) (*model.Currency, error)
}

type currencyRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewCurrencyRepo(db *db.DatabaseConnection) CurrencyRepo {
	return &currencyRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *currencyRepo) GetDefaultCurrency(ctx context.Context) (*model.Currency, error) {
	currency := new(model.Currency)

	q := `SELECT code, name, minor_unit, is_active, is_default
		FROM currency
		WHERE is_default = TRUE
		LIMIT 1`

	if err := r.roDb.GetContext(ctx, currency, q); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get default currency")
	}

	return currency, nil
}
