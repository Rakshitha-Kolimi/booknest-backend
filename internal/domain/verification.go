package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type VerificationTokenType string

const (
	VerificationEmail VerificationTokenType = "EMAIL_VERIFICATION"
	PasswordReset     VerificationTokenType = "PASSWORD_RESET"
	LoginOTP          VerificationTokenType = "LOGIN_OTP"
)

type VerificationToken struct {
	ID        uuid.UUID             `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID             `gorm:"type:uuid;index"`
	Type      VerificationTokenType `gorm:"type:varchar(32);index"`
	TokenHash string                `gorm:"not null"`
	ExpiresAt time.Time             `gorm:"index"`
	IsUsed    bool                  `gorm:"default:false"`

	UsedAt   *time.Time
	Metadata map[string]string `gorm:"type:jsonb"`

	BaseEntity
}

type VerificationTokenRepository interface {
	FindByUserIDAndType(
		ctx context.Context,
		userID uuid.UUID,
		tokenType VerificationTokenType,
	) (*VerificationToken, error)

	Create(ctx context.Context, token *VerificationToken) error
	Update(ctx context.Context, token *VerificationToken) error
	Delete(ctx context.Context, id uuid.UUID) error
}
