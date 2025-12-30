package repositories

import (
"context"
"time"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
"project-management/models"
)

type SessionRepository interface {
Create(ctx context.Context, session *models.Session) error
GetByRefreshToken(ctx context.Context, tokenHash string) (*models.Session, error)
Revoke(ctx context.Context, tokenHash string) error
DeleteExpired(ctx context.Context) error
}

type sessionRepository struct {
db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
query := `
INSERT INTO sessions (id, user_id, refresh_token_hash, user_agent, ip_address, created_at, expires_at, revoked)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at
`

session.ID = uuid.New()
session.CreatedAt = time.Now()

return r.db.QueryRow(ctx, query,
session.ID,
session.UserID,
session.RefreshTokenHash,
session.UserAgent,
session.IPAddress,
session.CreatedAt,
session.ExpiresAt,
session.Revoked,
).Scan(&session.ID, &session.CreatedAt)
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, tokenHash string) (*models.Session, error) {
query := `
SELECT id, user_id, refresh_token_hash, user_agent, ip_address, created_at, expires_at, revoked
FROM sessions
WHERE refresh_token_hash = $1
`

session := &models.Session{}
err := r.db.QueryRow(ctx, query, tokenHash).Scan(
&session.ID,
&session.UserID,
&session.RefreshTokenHash,
&session.UserAgent,
&session.IPAddress,
&session.CreatedAt,
&session.ExpiresAt,
&session.Revoked,
)
if err != nil {
return nil, err
}

return session, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, tokenHash string) error {
query := `
UPDATE sessions
SET revoked = true
WHERE refresh_token_hash = $1
`

_, err := r.db.Exec(ctx, query, tokenHash)
return err
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
query := `
DELETE FROM sessions
WHERE expires_at < $1 OR revoked = true
`

_, err := r.db.Exec(ctx, query, time.Now())
return err
}
