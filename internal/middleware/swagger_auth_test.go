package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerAuthMiddleware(t *testing.T) {
	// Gin test mode
	gin.SetMode(gin.TestMode)

	// Set env vars
	t.Setenv("SWAGGER_USER", "admin")
	t.Setenv("SWAGGER_PASSWORD", "secret")

	router := gin.New()
	router.GET(
		"/swagger/*any",
		SwaggerAuthMiddleware(),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		},
	)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid credentials",
			authHeader:     "Basic YWRtaW46c2VjcmV0", // base64(admin:secret)
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid password",
			authHeader:     "Basic YWRtaW46d3Jvbmc=", // admin:wrong
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/swagger/index.html",
				nil,
			)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
