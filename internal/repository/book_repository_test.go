package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

// ------------------ BOOK TESTS ------------------

func TestGetBook(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "books")
	seed := seedBook(t)

	repo := NewBookRepositoryImpl(testDB)
	ctx := context.Background()

	got, err := repo.GetBook(ctx, seed.ID)
	require.NoError(t, err)
	require.Equal(t, seed.Title, got.Title)

	defer cleanTable(t, "books")
}

func TestGetBooks(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "books")
	seed := seedBook(t)

	repo := NewBookRepositoryImpl(testDB)
	ctx := context.Background()

	books, err := repo.GetBooks(ctx)
	require.NoError(t, err)
	require.Len(t, books, 1)
	require.Equal(t, seed.ID, books[0].ID)
}

func TestCreateBook(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "books")
	repo := NewBookRepositoryImpl(testDB)
	ctx := context.Background()

	book := &domain.Book{
		Title:  "New Book",
		Author: "John Doe",
		Price:  20.5,
		Stock:  5,
	}

	err := repo.CreateBook(ctx, book)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, book.ID)
}

func TestUpdateBook(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "books")
	seed := seedBook(t)

	repo := NewBookRepositoryImpl(testDB)
	ctx := context.Background()

	seed.Title = "Updated Title"
	err := repo.UpdateBook(ctx, &seed)
	require.NoError(t, err)

	updated, err := repo.GetBook(ctx, seed.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Title", updated.Title)
}

func TestDeleteBook(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "books")
	seed := seedBook(t)

	repo := NewBookRepositoryImpl(testDB)
	ctx := context.Background()

	err := repo.DeleteBook(ctx, seed.ID)
	require.NoError(t, err)

	_, err = repo.GetBook(ctx, seed.ID)
	require.Error(t, err) // Should fail because book is soft-deleted
}
