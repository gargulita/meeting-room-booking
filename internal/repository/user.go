package repository

import (
	"context"
	"database/sql"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, role, password_hash, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING created_at
	`
	return r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Role, user.PasswordHash).Scan(&user.CreatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, role, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	var user domain.User
	var passwordHash sql.NullString
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Role, &passwordHash, &user.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, role, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	var user domain.User
	var passwordHash sql.NullString
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Role, &passwordHash, &user.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	return &user, nil
}
