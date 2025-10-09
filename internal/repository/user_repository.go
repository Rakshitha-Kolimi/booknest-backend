package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepositoryImpl(db *pgxpool.Pool) domain.UserRepository {
	return &userRepositoryImpl{db: db}
}

// GetUsers returns all non-deleted users
func (r *userRepositoryImpl) GetUsers(ctx context.Context) ([]domain.User, error) {
	tx := getTx(ctx)

	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var rows pgx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Query(ctx, query)
	} else {
		rows, err = r.db.Query(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUser retrieves a user by ID
func (r *userRepositoryImpl) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	tx := getTx(ctx)

	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user domain.User
	var err error
	if tx != nil {
		err = tx.QueryRow(ctx, query, id).Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(ctx, query, id).Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.Email,
			&user.Password, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	return user, err
}

// CreateUser inserts a new user into the database
func (r *userRepositoryImpl) CreateUser(ctx context.Context, entity *domain.User) error {
	tx := getTx(ctx)

	query := `
		INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	args := []interface{}{entity.FirstName, entity.LastName, entity.Email, entity.Password}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
}

// UpdateUser updates user info
func (r *userRepositoryImpl) UpdateUser(ctx context.Context, entity *domain.User) error {
	tx := getTx(ctx)

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	args := []interface{}{
		entity.FirstName, entity.LastName, entity.Email, entity.Password, entity.ID,
	}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&entity.UpdatedAt)
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&entity.UpdatedAt)
}

// DeleteUser performs a soft delete
func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tx := getTx(ctx)

	query := `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{id}

	var cmdTag pgconn.CommandTag
	var err error

	if tx != nil {
		cmdTag, err = tx.Exec(ctx, query, args...)
	} else {
		cmdTag, err = r.db.Exec(ctx, query, args...)
	}
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
