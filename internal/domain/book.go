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
	DeleteBook(ctx context.Context, id uuid.UUID) (err error)
	GetBooks(ctx context.Context) ([]Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (Book, error)
	UpdateBook(ctx context.Context, entity *Book) (err error)
}

type BookService interface {
	CreateBook(ctx context.Context, in BookInput) (book Book, err error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	GetBooks(ctx context.Context) ([]Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (Book, error)
	UpdateBook(ctx context.Context, id uuid.UUID, entity BookInput) (book Book, err error)
}

type BookController interface {
	AddBook(ctx *gin.Context)
	GetBooks(ctx *gin.Context)
	GetBook(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
}
