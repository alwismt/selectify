package repo

import (
	"context"
	"database/sql"
	"errors"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRoleRepo interface {
	InsertUserRoleForCustomerWithTx(ctx context.Context, tx sqlx.ExecerContext, role *model.UserRole) error
	GetUserRoleByUserID(ctx context.Context, id uint) (user *model.UserRole, err error)
}

type userRoleRepo struct {
	roDb *sqlx.DB
	rwDB *sqlx.DB
}

func NewUserRoleRepo(db *db.DatabaseConnection) UserRoleRepo {
	return &userRoleRepo{
		roDb: db.RoDb,
		rwDB: db.RwDb,
	}
}

func (r *userRoleRepo) InsertUserRoleForCustomerWithTx(ctx context.Context, tx sqlx.ExecerContext, role *model.UserRole) error {
	query := `INSERT INTO user_role(user_id, role) VALUES ($1, $2)`
	_, err := tx.ExecContext(ctx, query, role.UserID, role.Role)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to insert user role for user %d", role.UserID)
	}
	return nil
}

func (r *userRoleRepo) GetUserRoleByUserID(ctx context.Context, id uint) (user *model.UserRole, err error) {
	user = new(model.UserRole)
	err = r.rwDB.QueryRowxContext(ctx, "SELECT * FROM user_role WHERE user_id = $1", id).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = logger.Errorf(ctx, err, "User role not found for id %d", id)
			return nil, httpx.ErrUserNotFound
		}
		err = logger.Errorf(ctx, err, "Failed to get user role by user id %d", id)
		return
	}
	return
}
