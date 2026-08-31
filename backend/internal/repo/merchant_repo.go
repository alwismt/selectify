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
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error
	UpdateMerchantWithTx(ctx context.Context, merchant *model.Merchant, tx sqlx.QueryerContext) error
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

func (r *merchantRepo) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	return r.updateMerchantWithTx(ctx, merchant, r.rwDb)
}

func (r *merchantRepo) UpdateMerchantWithTx(ctx context.Context, merchant *model.Merchant, tx sqlx.QueryerContext) error {
	return r.updateMerchantWithTx(ctx, merchant, tx)
}

func (r *merchantRepo) updateMerchantWithTx(ctx context.Context, merchant *model.Merchant, tx sqlx.QueryerContext) error {
	q := `UPDATE merchant SET
			name = $1,
			description = $2,
			country_code = $3,
			logo = $4,
			updated_at = NOW()
		WHERE merchant_id = $5
		RETURNING
			updated_at`
	err := tx.QueryRowxContext(ctx, q, merchant.Name, merchant.Description, merchant.CountryCode, merchant.Logo, merchant.MerchantID).StructScan(merchant)
	if err != nil {
		return logger.Errorf(ctx, err, "Failed to update merchant %d", merchant.MerchantID)
	}
	return nil
}
