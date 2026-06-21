package auth

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByEmail(ctx context.Context, email string) (User, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	return err
}

func (r *postgresRepository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	return u, err
}

var ErrUserNotFound = errors.New("usuário não encontrado")