package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type bookService struct {
	repo domain.BookRepository
	db   *gorm.DB
}

func NewBookService(repo domain.BookRepository, db *gorm.DB) domain.BookService {
	return &bookService{
		repo: repo,
		db:   db,
	}
}

func (s *bookService) CreateBook(ctx context.Context, input domain.BookInput) (*domain.Book, error) {
	book := &domain.Book{
		ID:                 uuid.New(),
		Name:               input.Name,
		AuthorName:         input.AuthorName,
		AvailableStock:     input.AvailableStock,
		ImageURL:           input.ImageURL,
		IsActive:           input.IsActive,
		Description:        input.Description,
		ISBN:               input.ISBN,
		Price:              input.Price,
		DiscountPercentage: input.DiscountPercentage,
		PublisherID:        input.PublisherID,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(book).Error; err != nil {
			return err
		}

		if len(input.CategoryIDs) > 0 {
			for _, cid := range input.CategoryIDs {
				bc := domain.BookCategory{
					BookID:     book.ID,
					CategoryID: cid,
				}
				if err := tx.Create(&bc).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) GetBook(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *bookService) ListBooks(ctx context.Context, limit, offset int) ([]domain.Book, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *bookService) FilterByCriteria(
	ctx context.Context,
	filter domain.BookFilter,
	q domain.QueryOptions,
) (*domain.BookSearchResult, error) {

	books, total, err := s.repo.FilterByCriteria(ctx, filter, q)
	if err != nil {
		return nil, err
	}

	return &domain.BookSearchResult{
		Items:  books,
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}, nil
}
