package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName      string     `gorm:"not null" json:"first_name"`
	LastName       string     `gorm:"not null" json:"last_name"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	Mobile         string     `gorm:"uniqueIndex;not null" json:"mobile"`
	Password       string     `gorm:"not null" json:"-"`
	LastLogin      *time.Time `json:"last_login,omitempty"`
	Role           UserRole   `gorm:"type:user_role;default:USER" json:"role"`
	IsActive       bool       `gorm:"default:false" json:"is_active"`
	EmailVerified  bool       `gorm:"default:false" json:"email_verified"`
	MobileVerified bool       `gorm:"default:false" json:"mobile_verified"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      time.Time  `json:"deleted_at"`
} // @name User

// UserInput is used for creating or updating users
type UserInput struct {
	FirstName string   `json:"first_name" binding:"required,min=3"`
	LastName  string   `json:"last_name" binding:"required"`
	Email     string   `json:"email" binding:"required,email"`
	Mobile    string   `json:"mobile" binding:"required,mobile"`
	Password  string   `json:"password" binding:"required,min=6"`
	Role      UserRole `json:"role" binding:"required"`
}

// ForgotPasswordInput is used for forgot password
type ForgotPasswordInput struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// LoginInput is used for login
type LoginInput struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password" binding:"required"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByMobile(ctx context.Context, mobile string) (User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
type UserService interface {
	Register(in UserInput) error
	VerifyEmail()
	VerifyMobile()
	Login(in LoginInput) (token string, err error)
	ForgotPassword(in ForgotPasswordInput) (err error)
	ResetPassword() error
}

type UserController interface {
	RegisterRoutes(r *gin.Engine)
}
