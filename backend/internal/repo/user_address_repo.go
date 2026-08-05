package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type UserAddressRepo interface {
	GetDefaultByUserID(ctx context.Context, userID uint) (*model.UserAddress, error)
	UpsertDefaultAddress(ctx context.Context, addr *model.UserAddress) error
}

type userAddressRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewUserAddressRepo(dbConn *db.DatabaseConnection) UserAddressRepo {
	return &userAddressRepo{
		rwDb: dbConn.RwDb,
		roDb: dbConn.RoDb,
	}
}

// GetDefaultByUserID returns the default address, or the most recently updated one if none is default.
func (r *userAddressRepo) GetDefaultByUserID(ctx context.Context, userID uint) (*model.UserAddress, error) {
	var addr model.UserAddress
	const q = `
		SELECT *
		FROM user_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, updated_at DESC, id DESC
		LIMIT 1;`

	err := r.roDb.GetContext(ctx, &addr, q, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, logger.Error(ctx, err, "failed to get user address")
	}
	return &addr, nil
}

// UpsertDefaultAddress updates the existing default (or latest) address, or inserts a new default.
func (r *userAddressRepo) UpsertDefaultAddress(ctx context.Context, addr *model.UserAddress) error {
	var existing model.UserAddress
	const findQ = `
		SELECT *
		FROM user_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, updated_at DESC, id DESC
		LIMIT 1;`

	err := r.rwDb.GetContext(ctx, &existing, findQ, addr.UserID)
	addr.IsDefault = true

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return logger.Error(ctx, err, "failed to find user address for upsert")
	}

	if err == nil {
		const q = `
			UPDATE user_addresses SET
				phone = :phone,
				line1 = :line1,
				line2 = :line2,
				city = :city,
				region = :region,
				postal_code = :postal_code,
				country_code = :country_code,
				is_default = true,
				updated_at = now()
			WHERE id = :id
			RETURNING id, created_at, updated_at;`

		addr.ID = existing.ID
		stmt, prepErr := r.rwDb.PrepareNamedContext(ctx, q)
		if prepErr != nil {
			return logger.Error(ctx, prepErr, "failed to prepare user address update")
		}
		defer func() {
			if closeErr := stmt.Close(); closeErr != nil {
				_ = logger.Error(ctx, closeErr, "failed to close statement")
			}
		}()

		if err = stmt.GetContext(ctx, addr, addr); err != nil {
			return logger.Error(ctx, err, "failed to update user address")
		}
		return nil
	}

	const insertQ = `
		INSERT INTO user_addresses (
			user_id, phone, line1, line2, city, region, postal_code, country_code, is_default
		) VALUES (
			:user_id, :phone, :line1, :line2, :city, :region, :postal_code, :country_code, :is_default
		)
		RETURNING id, created_at, updated_at;`

	stmt, prepErr := r.rwDb.PrepareNamedContext(ctx, insertQ)
	if prepErr != nil {
		return logger.Error(ctx, prepErr, "failed to prepare user address insert")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			_ = logger.Error(ctx, closeErr, "failed to close statement")
		}
	}()

	if err = stmt.GetContext(ctx, addr, addr); err != nil {
		return logger.Error(ctx, err, "failed to insert user address")
	}
	return nil
}
