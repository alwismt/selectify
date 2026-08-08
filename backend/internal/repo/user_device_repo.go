package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type UserDeviceRepo interface {
	Insert(ctx context.Context, d *model.UserDevice) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.UserDevice, error)
	TouchIfStale(ctx context.Context, deviceID uuid.UUID, userAgent, ip string, activityThrottle time.Duration) error
}

type userDeviceRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewUserDeviceRepo(db *db.DatabaseConnection) UserDeviceRepo {
	return &userDeviceRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *userDeviceRepo) Insert(ctx context.Context, d *model.UserDevice) error {
	query := `
		INSERT INTO user_device (
			device_id, user_id, device_token_hash, user_agent, last_ip, first_seen_at, last_seen_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.rwDb.ExecContext(ctx, query,
		d.DeviceId,
		d.UserId,
		d.DeviceTokenHash,
		d.UserAgent,
		d.LastIP,
		d.FirstSeenAt,
		d.LastSeenAt,
	)
	if err != nil {
		return logger.Errorf(ctx, err, "error inserting user device for user id: %d", d.UserId)
	}
	return nil
}

func (r *userDeviceRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*model.UserDevice, error) {
	device := new(model.UserDevice)
	query := `
		SELECT *
		FROM user_device
		WHERE device_token_hash = $1
		LIMIT 1
	`
	err := r.rwDb.GetContext(ctx, device, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Errorf(ctx, err, "error getting user device by token hash")
	}
	return device, nil
}

func (r *userDeviceRepo) TouchIfStale(ctx context.Context, deviceID uuid.UUID, userAgent, ip string, activityThrottle time.Duration) error {
	now := time.Now().UTC()
	staleBefore := now.Add(-activityThrottle)
	query := `
		UPDATE user_device
		SET last_seen_at = $1, user_agent = $2, last_ip = $3
		WHERE device_id = $4 AND last_seen_at <= $5
	`
	_, err := r.rwDb.ExecContext(ctx, query, now, userAgent, ip, deviceID, staleBefore)
	if err != nil {
		return logger.Errorf(ctx, err, "error touching user device %s", deviceID)
	}
	return nil
}
