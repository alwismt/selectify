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

type UserSessionRepo interface {
	InsertUserSession(ctx context.Context, s *model.UserSession) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.UserSession, error)
	RenewIfStale(ctx context.Context, sessionID uuid.UUID, idleTTL, activityThrottle time.Duration) (*model.UserSession, error)
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
		INSERT INTO user_session (
			user_id, session_id, session_token_hash, expires_at, absolute_expires_at,
			last_activity_at, remember_me, device_id, user_agent, ip_address
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.rwDb.ExecContext(ctx, query,
		s.UserId,
		s.SessionId,
		s.SessionTokenHash,
		s.ExpiresAt,
		s.AbsoluteExpiresAt,
		s.LastActivityAt,
		s.RememberMe,
		s.DeviceId,
		s.UserAgent,
		s.IpAddress,
	)
	if err != nil {
		return logger.Errorf(ctx, err, "error adding user session for user id: %d", s.UserId)
	}
	return nil
}

func (r *userSessionRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*model.UserSession, error) {
	userSession := new(model.UserSession)
	query := `
		SELECT user_session.*
		FROM user_session
		WHERE session_token_hash = $1 AND revoked_at IS NULL
		LIMIT 1
	`

	err := r.rwDb.GetContext(ctx, userSession, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.Errorf(ctx, err, "error getting user session for token hash")
		}
		return nil, logger.Errorf(ctx, err, "error getting user session for token hash")
	}

	return userSession, nil
}

// RenewIfStale extends idle expiry when last_activity_at is older than activityThrottle.
// Remember-me sessions are not idle-extended. Concurrent callers are serialized by the
// conditional UPDATE so only one write succeeds per throttle window.
func (r *userSessionRepo) RenewIfStale(ctx context.Context, sessionID uuid.UUID, idleTTL, activityThrottle time.Duration) (*model.UserSession, error) {
	now := time.Now().UTC()
	staleBefore := now.Add(-activityThrottle)
	candidateExpires := now.Add(idleTTL)

	query := `
		UPDATE user_session
		SET
			last_activity_at = $1,
			expires_at = LEAST($2::timestamptz, absolute_expires_at)
		WHERE session_id = $3
		  AND revoked_at IS NULL
		  AND remember_me = FALSE
		  AND last_activity_at <= $4
		  AND absolute_expires_at > $1
		RETURNING *
	`

	updated := new(model.UserSession)
	err := r.rwDb.GetContext(ctx, updated, query, now, candidateExpires, sessionID, staleBefore)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Errorf(ctx, err, "error renewing user session %s", sessionID)
	}
	return updated, nil
}

func (r *userSessionRepo) RevokeSession(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE user_session SET revoked_at = $1 WHERE session_id = $2`
	_, err := r.rwDb.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return logger.Errorf(ctx, err, "error revoking user session for session id: %s", id)
	}
	return nil
}
