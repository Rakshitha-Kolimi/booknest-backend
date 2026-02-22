package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupServer_Success(t *testing.T) {
	t.Setenv("SWAGGER_USER", "swagger")
	t.Setenv("SWAGGER_PASSWORD", "swagger-pass")

	originalConnectGORM := connectGORM
	t.Cleanup(func() {
		connectGORM = originalConnectGORM
	})

	connectGORM = func() (*gorm.DB, error) {
		return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	}

	router, err := SetupServer(&pgxpool.Pool{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestSetupServer_ConnectGORMError(t *testing.T) {
	originalConnectGORM := connectGORM
	t.Cleanup(func() {
		connectGORM = originalConnectGORM
	})

	connectGORM = func() (*gorm.DB, error) {
		return nil, errors.New("db unavailable")
	}

	_, err := SetupServer(&pgxpool.Pool{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
