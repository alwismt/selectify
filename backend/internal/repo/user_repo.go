package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/httpx"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type UserRepo interface {
	AddUserWithTx(ctx context.Context, tx sqlx.QueryerContext, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (user *model.User, err error)
	GetUserById(ctx context.Context, id uint) (user *model.User, err error)
	UpdatePasswordHash(ctx context.Context, userID uint, passwordHash string) error
}

type userRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewUserRepo(dbConn *db.DatabaseConnection) UserRepo {
	return &userRepo{
		rwDb: dbConn.RwDb,
		roDb: dbConn.RoDb,
	}
}

func (uR *userRepo) AddUserWithTx(ctx context.Context, tx sqlx.QueryerContext, user *model.User) error {
	query := `INSERT INTO users (email, password_hash, first_name, last_name, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`

	err := tx.QueryRowxContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Phone, user.Status).
		StructScan(user)

	if err != nil {
		// Postgres unique constraint
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				_ = logger.Error(ctx, err, "User already exists")
				return httpx.ErrUserAlreadyExists
			}
		}

		_ = logger.Error(ctx, err, "Failed to insert user")
		return err
	}

	return nil
}

func (uR *userRepo) GetUserByEmail(ctx context.Context, email string) (user *model.User, err error) {
	user = new(model.User)
	err = uR.roDb.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1", email).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = logger.Errorf(ctx, err, "User not found for email %s", email)
			return nil, httpx.ErrUserNotFound
		}
		err = logger.Errorf(ctx, err, "Failed to get user by email %s", email)
		return
	}
	return
}

func (uR *userRepo) GetUserById(ctx context.Context, id uint) (user *model.User, err error) {
	user = new(model.User)
	err = uR.rwDb.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", id).StructScan(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_ = logger.Errorf(ctx, err, "User not found for id %d", id)
			return nil, httpx.ErrUserNotFound
		}
		err = logger.Errorf(ctx, err, "Failed to get user by id %d", id)
		return
	}
	return
}

func (uR *userRepo) UpdatePasswordHash(ctx context.Context, userID uint, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	result, err := uR.rwDb.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return logger.Errorf(ctx, err, "Failed to update password hash for user %d", userID)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return logger.Errorf(ctx, err, "Failed to read rows affected updating password for user %d", userID)
	}
	if rows == 0 {
		return httpx.ErrUserNotFound
	}
	return nil
}

func (uR *userRepo) AddUserWithRole(ctx context.Context, tx sqlx.QueryerContext, user *model.User, role *model.UserRole) error {
	query := `WITH inserted_user AS (
			INSERT INTO users (email, password_hash, first_name, last_name, phone, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		)
		INSERT INTO user_role (user_id, role)
		SELECT id, $7 FROM inserted_user
		RETURNING user_id
	`

	var userID uint
	err := tx.QueryRowxContext(ctx, query, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Phone, user.Status, role.Role).Scan(&userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return httpx.ErrUserAlreadyExists
		}
		return err
	}
	user.ID = userID
	role.UserID = userID
	return nil
}
