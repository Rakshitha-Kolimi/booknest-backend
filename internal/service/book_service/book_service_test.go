package book_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type mockBookRepository struct {
	findByIDFunc         func(ctx context.Context, id uuid.UUID) (*domain.Book, error)
	listFunc             func(ctx context.Context, limit, offset int) ([]domain.Book, error)
	filterByCriteriaFunc func(ctx context.Context, filter domain.BookFilter, pagination domain.QueryOptions) ([]domain.Book, int64, error)
	deleteFunc           func(ctx context.Context, id uuid.UUID) error
}

func (m *mockBookRepository) Create(ctx context.Context, book *domain.Book) error {
	return nil
}

func (m *mockBookRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockBookRepository) List(ctx context.Context, limit, offset int) ([]domain.Book, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return []domain.Book{}, nil
}

func (m *mockBookRepository) FilterByCriteria(ctx context.Context, filter domain.BookFilter, pagination domain.QueryOptions) ([]domain.Book, int64, error) {
	if m.filterByCriteriaFunc != nil {
		return m.filterByCriteriaFunc(ctx, filter, pagination)
	}
	return []domain.Book{}, 0, nil
}

func (m *mockBookRepository) Update(ctx context.Context, book *domain.Book) error {
	return nil
}

func (m *mockBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestUniqueCategoryIDsRemovesNilAndDuplicates(t *testing.T) {
	first := uuid.New()
	second := uuid.New()

	result := uniqueCategoryIDs([]uuid.UUID{uuid.Nil, first, second, first, uuid.Nil, second})
	if len(result) != 2 {
		t.Fatalf("expected 2 unique category IDs, got %d", len(result))
	}
	if result[0] != first || result[1] != second {
		t.Fatalf("unexpected order/content: %+v", result)
	}

	empty := uniqueCategoryIDs(nil)
	if empty != nil {
		t.Fatalf("expected nil for empty input, got %+v", empty)
	}
}

func TestValidateCategories(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&domain.Category{}); err != nil {
		t.Fatalf("failed migration: %v", err)
	}

	idOne := uuid.New()
	idTwo := uuid.New()
	if err := db.Create(&domain.Category{ID: idOne, Name: "A"}).Error; err != nil {
		t.Fatalf("failed to seed category A: %v", err)
	}
	if err := db.Create(&domain.Category{ID: idTwo, Name: "B"}).Error; err != nil {
		t.Fatalf("failed to seed category B: %v", err)
	}

	if err := validateCategories(db, []uuid.UUID{idOne, idTwo}); err != nil {
		t.Fatalf("expected categories to validate, got %v", err)
	}

	missingErr := validateCategories(db, []uuid.UUID{idOne, uuid.New()})
	if missingErr == nil || missingErr.Error() != "one or more categories are invalid" {
		t.Fatalf("expected invalid categories error, got %v", missingErr)
	}
}

func TestBookServiceReadAndFilterPassThrough(t *testing.T) {
	bookID := uuid.New()
	expected := &domain.Book{ID: bookID, Name: "Book"}
	filter := domain.BookFilter{}
	query := domain.QueryOptions{Limit: 10, Offset: 5}

	repo := &mockBookRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Book, error) {
			if id != bookID {
				t.Fatalf("unexpected ID: %s", id)
			}
			return expected, nil
		},
		listFunc: func(ctx context.Context, limit, offset int) ([]domain.Book, error) {
			if limit != 10 || offset != 20 {
				t.Fatalf("unexpected pagination: %d/%d", limit, offset)
			}
			return []domain.Book{*expected}, nil
		},
		filterByCriteriaFunc: func(ctx context.Context, gotFilter domain.BookFilter, pagination domain.QueryOptions) ([]domain.Book, int64, error) {
			if pagination.Limit != query.Limit || pagination.Offset != query.Offset {
				t.Fatalf("unexpected query options: %+v", pagination)
			}
			return []domain.Book{*expected}, 1, nil
		},
		deleteFunc: func(ctx context.Context, id uuid.UUID) error {
			if id != bookID {
				t.Fatalf("unexpected delete ID: %s", id)
			}
			return nil
		},
	}

	svc := NewBookService(repo, nil)

	book, err := svc.GetBook(context.Background(), bookID)
	if err != nil || book.ID != bookID {
		t.Fatalf("unexpected GetBook result: %+v, err=%v", book, err)
	}

	books, err := svc.ListBooks(context.Background(), 10, 20)
	if err != nil || len(books) != 1 {
		t.Fatalf("unexpected ListBooks result: %+v, err=%v", books, err)
	}

	result, err := svc.FilterByCriteria(context.Background(), filter, query)
	if err != nil {
		t.Fatalf("unexpected filter error: %v", err)
	}
	if result.Total != 1 || result.Limit != query.Limit || result.Offset != query.Offset || len(result.Items) != 1 {
		t.Fatalf("unexpected filter result: %+v", result)
	}

	if err := svc.DeleteBook(context.Background(), bookID); err != nil {
		t.Fatalf("unexpected DeleteBook error: %v", err)
	}
}
