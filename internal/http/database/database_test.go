package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestConnect_Success(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_NAME", "booknest_test")
	os.Setenv("DB_PORT", "5432")

	// mock function
	original := newPgxPool
	defer func() { newPgxPool = original }()

	newPgxPool = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return &pgxpool.Pool{}, nil
	}

	pool, err := Connect()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if pool == nil {
		t.Fatalf("expected non-nil pool")
	}
}

func TestConnect_Fail_NewPool(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_PORT", "5432")

	original := newPgxPool
	defer func() { newPgxPool = original }()

	newPgxPool = func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, errors.New("mock connection error")
	}

	_, err := Connect()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestConnect_Fail_ParseConfig(t *testing.T) {
	os.Unsetenv("DB_HOST") // invalid DSN should trigger parse error
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PORT")

	_, err := Connect()
	if err == nil {
		t.Fatalf("expected parse config error, got nil")
	}
}