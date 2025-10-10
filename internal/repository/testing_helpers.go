package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

var testDB *pgxpool.Pool

// Initialize testDB safely for any test
func initTestDB(t *testing.T) *pgxpool.Pool {
	if testDB != nil {
		return testDB
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/booknest_test?sslmode=disable"
	}

	var err error
	testDB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect test DB: %v", err)
	}
	log.Println("✅ Connected to test DB")
	return testDB
}

// Clean table data between tests
func cleanTable(t *testing.T, tableName string) {
	if testDB == nil {
		t.Fatal("testDB is nil! Did you forget to call initTestDB(t)?")
	}

	_, err := testDB.Exec(context.Background(), "TRUNCATE "+tableName+" RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to truncate %s: %v", tableName, err)
	}
}

// Seed helpers
func seedBook(t *testing.T) domain.Book {
	initTestDB(t)
	ctx := context.Background()
	book := domain.Book{
		ID:        uuid.New(),
		Title:     "Seeded Test Book",
		Author:    "Jane Doe",
		Price:     15.50,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := testDB.Exec(ctx, `
		INSERT INTO books (id, title, author, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, book.ID, book.Title, book.Author, book.Price, book.Stock, book.CreatedAt, book.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to seed book: %v", err)
	}
	return book
}

func seedUser(t *testing.T) domain.User {
	initTestDB(t)
	ctx := context.Background()
	user := domain.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := testDB.Exec(ctx, `
		INSERT INTO users (id, first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return user
}
