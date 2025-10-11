package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"booknest/internal/domain"
)

type mockBookRepo struct {
	books       []domain.Book
	createErr   error
	getErr      error
	getBooksErr error
	updateErr   error
	deleteErr   error

	createdBook *domain.Book
	updatedBook *domain.Book
	deletedID   uuid.UUID
}

func (m *mockBookRepo) CreateBook(ctx context.Context, book *domain.Book) error {
	if m.createErr != nil {
		return m.createErr
	}
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
	}
	if book.CreatedAt.IsZero() {
		book.CreatedAt = time.Now()
	}
	m.createdBook = book
	m.books = append(m.books, *book)
	return nil
}

func (m *mockBookRepo) DeleteBook(ctx context.Context, id uuid.UUID) error {
	m.deletedID = id
	return m.deleteErr
}

func (m *mockBookRepo) GetBooks(ctx context.Context) ([]domain.Book, error) {
	if m.getBooksErr != nil {
		return nil, m.getBooksErr
	}
	return m.books, nil
}

func (m *mockBookRepo) GetBook(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	if m.getErr != nil {
		return domain.Book{}, m.getErr
	}
	for _, b := range m.books {
		if b.ID == id {
			return b, nil
		}
	}
	return domain.Book{}, errors.New("book not found")
}

func (m *mockBookRepo) UpdateBook(ctx context.Context, entity *domain.Book) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedBook = entity
	for i := range m.books {
		if m.books[i].ID == entity.ID {
			m.books[i] = *entity
			break
		}
	}
	return nil
}

func TestBookService_CRUD(t *testing.T) {
	repo := &mockBookRepo{books: []domain.Book{}}
	svc := NewBookServiceImpl(repo)

	in := domain.BookInput{Title: "Go Programming", Author: "R", Price: 29.99, Stock: 10}

	// Create
	created, err := svc.CreateBook(context.Background(), in)
	if err != nil {
		t.Fatalf("CreateBook error: %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatalf("expected created book to have ID assigned")
	}

	// GetBooks
	list, err := svc.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("GetBooks error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 book, got %d", len(list))
	}

	// GetBook
	got, err := svc.GetBook(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetBook error: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected id %s got %s", created.ID, got.ID)
	}

	// Update
	up := domain.BookInput{Title: "Advanced Go", Author: "R2", Price: 39.99, Stock: 5}
	updated, err := svc.UpdateBook(context.Background(), created.ID, up)
	if err != nil {
		t.Fatalf("UpdateBook error: %v", err)
	}
	if updated.Title != up.Title || updated.Author != up.Author {
		t.Fatalf("update didn't apply: %+v", updated)
	}

	// Delete
	if err := svc.DeleteBook(context.Background(), created.ID); err != nil {
		t.Fatalf("DeleteBook error: %v", err)
	}
	if repo.deletedID != created.ID {
		t.Fatalf("expected deleted id %s, got %s", created.ID, repo.deletedID)
	}
}
