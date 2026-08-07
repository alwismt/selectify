package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type PasswordResetRepo interface {
	Insert(ctx context.Context, reset *model.PasswordReset) error
	GetValidByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordReset, error)
	MarkUsed(ctx context.Context, passwordResetID string, usedAt time.Time) error
	InvalidateUnusedForUser(ctx context.Context, userID uint) error
}

type passwordResetRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewPasswordResetRepo(dbConn *db.DatabaseConnection) PasswordResetRepo {
	return &passwordResetRepo{
		rwDb: dbConn.RwDb,
		roDb: dbConn.RoDb,
	}
}

func (r *passwordResetRepo) Insert(ctx context.Context, reset *model.PasswordReset) error {
	query := `
		INSERT INTO password_reset (password_reset_id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.rwDb.ExecContext(ctx, query, reset.PasswordResetID, reset.UserID, reset.TokenHash, reset.ExpiresAt)
	if err != nil {
		return logger.Errorf(ctx, err, "error inserting password reset for user id: %d", reset.UserID)
	}
	return nil
}

func (r *passwordResetRepo) GetValidByTokenHash(ctx context.Context, tokenHash string) (*model.PasswordReset, error) {
	reset := new(model.PasswordReset)
	query := `
		SELECT password_reset_id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset
		WHERE token_hash = $1
		  AND used_at IS NULL
		  AND expires_at > NOW()
		LIMIT 1
	`
	err := r.rwDb.GetContext(ctx, reset, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Errorf(ctx, err, "error getting password reset by token hash")
	}
	return reset, nil
}

func (r *passwordResetRepo) MarkUsed(ctx context.Context, passwordResetID string, usedAt time.Time) error {
	query := `UPDATE password_reset SET used_at = $1 WHERE password_reset_id = $2 AND used_at IS NULL`
	result, err := r.rwDb.ExecContext(ctx, query, usedAt, passwordResetID)
	if err != nil {
		return logger.Errorf(ctx, err, "error marking password reset used: %s", passwordResetID)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return logger.Errorf(ctx, err, "error reading rows affected for password reset: %s", passwordResetID)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *passwordResetRepo) InvalidateUnusedForUser(ctx context.Context, userID uint) error {
	query := `
		UPDATE password_reset
		SET expires_at = NOW()
		WHERE user_id = $1
		  AND used_at IS NULL
		  AND expires_at > NOW()
	`
	_, err := r.rwDb.ExecContext(ctx, query, userID)
	if err != nil {
		return logger.Errorf(ctx, err, "error invalidating unused password resets for user id: %d", userID)
	}
	return nil
}
