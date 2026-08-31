package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type CountryRepo interface {
	CountryExists(ctx context.Context, code string) bool
	Countries(ctx context.Context) (*model.Countries, error)
}

type countryRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewCountryRepo(db *db.DatabaseConnection) CountryRepo {
	return &countryRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *countryRepo) CountryExists(ctx context.Context, code string) bool {
	const q = `SELECT EXISTS (
			SELECT 1
			FROM country
			WHERE code = $1
			  AND is_active = TRUE)`

	var exists bool
	if err := r.roDb.GetContext(ctx, &exists, q, code); err != nil {
		_ = logger.Error(ctx, err, "failed to get country")
		return false
	}

	return exists
}

func (r *countryRepo) Countries(ctx context.Context) (*model.Countries, error) {
	countries := new(model.Countries)
	const q = `SELECT code, name, is_active, created_at, updated_at FROM country WHERE is_active = TRUE`
	if err := r.roDb.SelectContext(ctx, countries, q); err != nil {
		_ = logger.Error(ctx, err, "failed to get countries")
		return nil, err
	}
	return countries, nil
}
