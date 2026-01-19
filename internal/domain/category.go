package domain

import (
	"time"

	"github.com/google/uuid"
)

// Categoty defines model for Category
type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type CategoryRepository interface{}
type CategoryService interface{}
type CategoryController interface{}
