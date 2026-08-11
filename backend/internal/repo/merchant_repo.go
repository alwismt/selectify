package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type MerchantRepo interface {
	GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error)
}

type merchantRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewMerchantRepo(db *db.DatabaseConnection) MerchantRepo {
	return &merchantRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *merchantRepo) GetMerchant(ctx context.Context, merchantID uint) (*model.Merchant, error) {
	var merchant model.Merchant
	if err := r.rwDb.GetContext(ctx, &merchant, `SELECT * FROM merchant WHERE merchant_id = $1`, merchantID); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to get merchant by id %d", merchantID)
	}
	return &merchant, nil
}
