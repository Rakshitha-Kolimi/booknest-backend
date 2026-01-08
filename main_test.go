package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestSetupServer_Success(t *testing.T) {
	// mock environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_NAME", "booknest_test")
	os.Setenv("DB_PORT", "5432")

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
