package repo

import (
	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRoleRepo interface {
	InsertUserRoleForCustomerWithTx(ctx context.Context, tx sqlx.ExecerContext, role *model.UserRole) error
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

func (uR *userRoleRepo) InsertUserRoleForCustomerWithTx(ctx context.Context, tx sqlx.ExecerContext, role *model.UserRole) error {
	query := `INSERT INTO user_role(user_id, role) VALUES ($1, $2)`
	_, err := tx.ExecContext(ctx, query, role.UserID, role.Role)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to insert user role for user %d", role.UserID)
	}
	return nil
}
