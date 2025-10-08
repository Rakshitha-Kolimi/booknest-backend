package service

import (
	"context"

	"github.com/google/uuid"

	"booknest/internal/domain"
)

type bookServiceImpl struct {
	r domain.BookRepository
}

func NewBookServiceImpl(r domain.BookRepository) domain.BookService {
	return &bookServiceImpl{r: r}
}

// ✅ Use ctx from caller instead of context.TODO()
func (s *bookServiceImpl) GetBook(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	return s.r.GetBook(ctx, id)
}

func (s *bookServiceImpl) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return s.r.GetBooks(ctx)
}

func (s *bookServiceImpl) CreateBook(ctx context.Context, in domain.BookInput) (domain.Book, error) {
	book := domain.Book{
		Title:  in.Title,
		Author: in.Author,
		Price:  in.Price,
		Stock:  in.Stock,
	}

	if err := s.r.CreateBook(ctx, &book); err != nil {
		return book, err
	}
	return book, nil
}

func (s *bookServiceImpl) UpdateBook(ctx context.Context, id uuid.UUID, in domain.BookInput) (domain.Book, error) {
	book, err := s.r.GetBook(ctx, id)
	if err != nil {
		return book, err
	}

	book.Title = in.Title
	book.Author = in.Author
	book.Price = in.Price
	book.Stock = in.Stock

	if err := s.r.UpdateBook(ctx, &book); err != nil {
		return book, err
	}
	return book, nil
}

func (s *bookServiceImpl) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.r.DeleteBook(ctx, id)
}