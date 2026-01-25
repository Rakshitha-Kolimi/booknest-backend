package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func SwaggerAuthMiddleware() gin.HandlerFunc {
	user := os.Getenv("SWAGGER_USER")
	pass := os.Getenv("SWAGGER_PASSWORD")

	return gin.BasicAuth(gin.Accounts{
		user: pass,
	})
}
