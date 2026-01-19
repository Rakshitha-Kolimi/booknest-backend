package domain

import (
	"time"

	"github.com/google/uuid"
)

// Publisher defines model for Publisher
type Publisher struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LegalName   string    `gorm:"not null" json:"legal_name"`
	TradingName string    `gorm:"not null" json:"trading_name"`
	Email       string    `gorm:"not null" json:"email"`
	Mobile      string    `gorm:"not null" json:"mobile"`
	Address     string    `gorm:"type:text" json:"address"`
	City        string    `gorm:"not null" json:"city"`
	State       string    `gorm:"not null" json:"state"`
	Country     string    `gorm:"not null" json:"country"`
	Zipcode     string    `gorm:"not null" json:"zipcode"`
	IsActive    bool      `gorm:"default:false" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
} // @name Publisher

type PublisherRepository interface{}
type PublisherService interface{}
type PublisherController interface{}
