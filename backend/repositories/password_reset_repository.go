package repositories

import (
	"context"
	"project-management/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetByToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, tokenHash string) error
	DeleteExpired(ctx context.Context) (int, error)
}

type passwordResetRepository struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	query := "INSERT INTO password_reset_tokens (id, user_id, token_hash, created_at, expires_at, used) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at"
	token.ID = uuid.New()
	token.CreatedAt = time.Now()
	return r.db.QueryRow(ctx, query, token.ID, token.UserID, token.TokenHash, token.CreatedAt, token.ExpiresAt, token.Used).Scan(&token.ID, &token.CreatedAt)
}

func (r *passwordResetRepository) GetByToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	query := "SELECT id, user_id, token_hash, created_at, expires_at, used FROM password_reset_tokens WHERE token_hash = $1"
	token := &models.PasswordResetToken{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.CreatedAt, &token.ExpiresAt, &token.Used)
	return token, err
}

func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, tokenHash string) error {
	query := "UPDATE password_reset_tokens SET used = true WHERE token_hash = $1"
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

func (r *passwordResetRepository) DeleteExpired(ctx context.Context) (int, error) {
	query := "DELETE FROM password_reset_tokens WHERE expires_at < $1 OR used = true"
	result, err := r.db.Exec(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}

	return int(result.RowsAffected()), nil
}
