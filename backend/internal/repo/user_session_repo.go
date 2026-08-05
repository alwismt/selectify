package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type UserSessionRepo interface {
	InsertUserSession(ctx context.Context, role *model.UserSession) error
	GetBySessionId(ctx context.Context, id string) (*model.UserSession, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
}

type userSessionRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewUserSessionRepo(db *db.DatabaseConnection) UserSessionRepo {
	return &userSessionRepo{
		rwDb: db.RwDb,
		roDb: db.RoDb,
	}
}

func (r *userSessionRepo) InsertUserSession(ctx context.Context, s *model.UserSession) error {
	query := `
		INSERT INTO user_session (user_id, session_id, expires_at, user_agent, ip_address)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.rwDb.ExecContext(ctx, query, s.UserId, s.SessionId, s.ExpiresAt, s.UserAgent, s.IpAddress)
	if err != nil {
		return logger.Errorf(ctx, err, "error adding user session for user id: %d", s.UserId)
	}
	return err
}

func (r *userSessionRepo) GetBySessionId(ctx context.Context, id string) (*model.UserSession, error) {
	userSession := new(model.UserSession)
	query := `SELECT user_session.* FROM user_session WHERE session_id = $1 AND revoked_at IS NULL LIMIT 1`

	err := r.rwDb.GetContext(ctx, userSession, query, id)
	if err != nil {
		return nil, logger.Errorf(ctx, err, "error getting user session for session id: %s", id)
	}

	return userSession, nil
}

func (r *userSessionRepo) RevokeSession(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE user_session SET revoked_at = $1 WHERE session_id = $2`
	_, err := r.rwDb.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return logger.Errorf(ctx, err, "error revoking user session for session id: %s", id)
	}
	return nil
}
