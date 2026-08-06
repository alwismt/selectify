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

type UserFileRepo interface {
	AddUserFileWithTx(ctx context.Context, tx sqlx.ExecerContext, file *model.UserFile) error
	GetByUserID(ctx context.Context, userID uint) (*model.UserFile, error)
	DeleteByUserIDWithTx(ctx context.Context, tx sqlx.ExecerContext, userID uint) error
}

type userFileRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewUserFileRepo(dbConn *db.DatabaseConnection) UserFileRepo {
	return &userFileRepo{
		rwDb: dbConn.RwDb,
		roDb: dbConn.RoDb,
	}
}

func (r *userFileRepo) AddUserFileWithTx(ctx context.Context, tx sqlx.ExecerContext, file *model.UserFile) error {
	query := `INSERT INTO user_file (
			user_file_id,
			user_id,
			content_type
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			user_file_id = EXCLUDED.user_file_id,
			content_type = EXCLUDED.content_type,
			updated_at = NOW()`

	_, err := tx.ExecContext(ctx, query, file.ID, file.UserID, file.ContentType)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to upsert user file")
		return err
	}

	return nil
}

func (r *userFileRepo) GetByUserID(ctx context.Context, userID uint) (*model.UserFile, error) {
	userFile := new(model.UserFile)
	err := r.roDb.QueryRowxContext(ctx, "SELECT * FROM user_file WHERE user_id = $1", userID).StructScan(userFile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		_ = logger.Error(ctx, err, "Failed to get user file by user id")
		return nil, err
	}

	return userFile, nil
}

func (r *userFileRepo) DeleteByUserIDWithTx(ctx context.Context, tx sqlx.ExecerContext, userID uint) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM user_file WHERE user_id = $1", userID)
	if err != nil {
		_ = logger.Error(ctx, err, "Failed to delete user file")
		return err
	}

	return nil
}
