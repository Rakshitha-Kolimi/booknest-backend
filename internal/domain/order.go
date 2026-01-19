package domain

import (
	"time"

	"github.com/google/uuid"
)

// Order defines the model for Order
type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrderNumber   string         `gorm:"uniqueIndex;not null" json:"order_number"`
	TotalPrice    float64        `gorm:"type:numeric(10,2)" json:"total_price"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID"`
	PaymentMethod *PaymentMethod `gorm:"type:payment_method" json:"payment_method,omitempty"`
	PaymentStatus *PaymentStatus `gorm:"type:payment_status" json:"payment_status,omitempty"`
	Status        OrderStatus    `gorm:"type:order_status;default:PENDING" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     time.Time      `json:"deleted_at,omitempty"`
} // @name Order

// OrderItem defines model for OrderItem
type OrderItem struct {
	OrderID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"order_id"`
	BookID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"book_id"`
	PurchaseCount int       `gorm:"check:purchase_count > 0" json:"purchase_count"`
	PurchasePrice float64   `gorm:"type:numeric(10,2)" json:"purchase_price"`
	TotalPrice    float64   `gorm:"type:numeric(10,2)" json:"total_price"`
	Book          Book      `gorm:"foreignKey:BookID"`
	Order         Order     `gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at,omitempty"`
} // @name OrderItem

type OrderRepository interface{}
type OrderService interface{}
type OrderController interface{}
