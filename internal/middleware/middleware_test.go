package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test that LoggingMiddleware does not break request flow and returns handler response
func TestLoggingMiddleware_Basic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoggingMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Fatalf("expected body to contain pong, got %s", w.Body.String())
	}
}

// Test ErrorHandler returns JSON for errors accumulated in context
func TestErrorHandler_ReportsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.GET("/err", func(c *gin.Context) {
		c.Error(errors.New("bad request"))
	})

	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Fatalf("expected 400 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bad request") {
		t.Fatalf("expected error message in body, got %s", w.Body.String())
	}
}
