package repository

import (
	"context"
	"fmt"

	"cbt/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository handles user database operations.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database and returns the created user.
func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, created_at, updated_at
	`

	var user models.User
	user.Password = passwordHash // Store hash temporarily for return

	if err := r.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail fetches a user by email (includes password hash).
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	if err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

// GetUserByID fetches a user by ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}
