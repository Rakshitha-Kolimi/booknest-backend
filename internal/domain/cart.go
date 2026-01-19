package domain

import (
	"time"

	"github.com/google/uuid"
)

// Cart defines the model for Cart
type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} // @name Cart

// CartItem defines model for cart item
type CartItem struct {
	CartID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"cart_id"`
	BookID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	Count     int       `gorm:"check:count > 0" json:"count"`
	CartPrice float64   `gorm:"type:numeric(10,2)" json:"cart_price"`
	Book      Book      `gorm:"foreignKey:BookID"`
	Cart      Cart      `gorm:"foreignKey:CartID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
} // @name CartItem

type CartRepository interface{}
type CartService interface{}
type CartController interface{}
