package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// UserInput is used for creating or updating users
type UserInput struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

// Repository interface for data operations
type UserRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, entity *User) error
	UpdateUser(ctx context.Context, entity *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Service interface for business logic
type UserService interface {
	RegisterUser(ctx context.Context, in UserInput) (string, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
