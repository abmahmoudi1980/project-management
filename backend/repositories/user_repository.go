package repositories

import (
	"context"
	"fmt"
	"time"

	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, limit, offset int, role string, isActive *bool) ([]*models.User, int, error)
	UpdateFailedAttempts(ctx context.Context, userID uuid.UUID, attempts int) error
	LockAccount(ctx context.Context, userID uuid.UUID, lockUntil time.Time) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users (id, username, email, password_hash, role, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at"

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return r.db.QueryRow(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := "SELECT id, username, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at, last_login_at FROM users WHERE id = $1"

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, username, email, password_hash, role, is_active, failed_login_attempts, locked_until, created_at, updated_at, last_login_at FROM users WHERE email = $1"

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := "UPDATE users SET username = $1, email = $2, role = $3, is_active = $4, password_hash = $5, last_login_at = $6, updated_at = $7 WHERE id = $8"

	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Role, user.IsActive, user.PasswordHash, user.LastLoginAt, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int, role string, isActive *bool) ([]*models.User, int, error) {
	query := "SELECT id, username, email, role, is_active, created_at, updated_at, last_login_at FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		countQuery += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, role)
		argCount++
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *isActive)
		argCount++
	}

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) UpdateFailedAttempts(ctx context.Context, userID uuid.UUID, attempts int) error {
	query := "UPDATE users SET failed_login_attempts = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.Exec(ctx, query, attempts, time.Now(), userID)
	return err
}

func (r *userRepository) LockAccount(ctx context.Context, userID uuid.UUID, lockUntil time.Time) error {
	query := "UPDATE users SET locked_until = $1, updated_at = $2 WHERE id = $3"
	_, err := r.db.Exec(ctx, query, lockUntil, time.Now(), userID)
	return err
}
