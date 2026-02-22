package book_service

import (
	"context"
	"fmt"

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
	authorID := uuid.Nil
	categoryIDs := uniqueCategoryIDs(input.CategoryIDs)

	book := &domain.Book{
		ID:   uuid.New(),
		Name: input.Name,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.AuthorID != nil && *input.AuthorID != uuid.Nil {
			authorID = *input.AuthorID
			var author domain.Author
			if err := tx.First(&author, "id = ?", authorID).Error; err != nil {
				return err
			}
		} else {
			var author domain.Author
			err := tx.Where("LOWER(name) = LOWER(?)", input.AuthorName).
				First(&author).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}

				author = domain.Author{
					ID:   uuid.New(),
					Name: input.AuthorName,
				}

				if err := tx.Create(&author).Error; err != nil {
					return err
				}
			}

			authorID = author.ID
		}

		book.AuthorID = authorID
		book.AvailableStock = input.AvailableStock
		book.ImageURL = input.ImageURL
		book.IsActive = input.IsActive
		book.Description = input.Description
		book.ISBN = input.ISBN
		book.Price = input.Price
		book.DiscountPercentage = input.DiscountPercentage
		book.PublisherID = input.PublisherID

		if err := tx.Create(book).Error; err != nil {
			return err
		}

		if len(categoryIDs) > 0 {
			if err := validateCategories(tx, categoryIDs); err != nil {
				return err
			}
			for _, cid := range categoryIDs {
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

func (s *bookService) UpdateBook(
	ctx context.Context,
	id uuid.UUID,
	input domain.BookInput,
) (*domain.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	authorID := book.AuthorID
	categoryIDs := uniqueCategoryIDs(input.CategoryIDs)

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.AuthorID != nil && *input.AuthorID != uuid.Nil {
			authorID = *input.AuthorID
			var author domain.Author
			if err := tx.First(&author, "id = ?", authorID).Error; err != nil {
				return err
			}
		} else {
			var author domain.Author
			err := tx.Where("LOWER(name) = LOWER(?)", input.AuthorName).
				First(&author).Error
			if err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}

				author = domain.Author{
					ID:   uuid.New(),
					Name: input.AuthorName,
				}
				if err := tx.Create(&author).Error; err != nil {
					return err
				}
			}

			authorID = author.ID
		}

		book.Name = input.Name
		book.AuthorID = authorID
		book.AvailableStock = input.AvailableStock
		book.ImageURL = input.ImageURL
		book.IsActive = input.IsActive
		book.Description = input.Description
		book.ISBN = input.ISBN
		book.Price = input.Price
		book.DiscountPercentage = input.DiscountPercentage
		book.PublisherID = input.PublisherID

		if err := tx.Save(book).Error; err != nil {
			return err
		}

		if err := tx.Where("book_id = ?", book.ID).Delete(&domain.BookCategory{}).Error; err != nil {
			return err
		}

		if len(categoryIDs) > 0 {
			if err := validateCategories(tx, categoryIDs); err != nil {
				return err
			}
			for _, cid := range categoryIDs {
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

func (s *bookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func uniqueCategoryIDs(categoryIDs []uuid.UUID) []uuid.UUID {
	if len(categoryIDs) == 0 {
		return nil
	}

	seen := make(map[uuid.UUID]struct{}, len(categoryIDs))
	unique := make([]uuid.UUID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique
}

func validateCategories(tx *gorm.DB, categoryIDs []uuid.UUID) error {
	var count int64
	if err := tx.Model(&domain.Category{}).Where("id IN ?", categoryIDs).Count(&count).Error; err != nil {
		return err
	}

	if int(count) != len(categoryIDs) {
		return fmt.Errorf("one or more categories are invalid")
	}

	return nil
}
