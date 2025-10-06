package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type BookInput struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

type BookRepository interface {
	CreateBook(ctx context.Context, book *Book) error
}

type BookService interface {
	CreateBook(in BookInput) (book Book, err error)
}

type BookController interface{
	AddBook(ctx *gin.Context)
}