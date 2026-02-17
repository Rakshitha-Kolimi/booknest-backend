package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"booknest/internal/domain"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		raw, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		role := domain.UserRole("")
		switch v := raw.(type) {
		case string:
			role = domain.UserRole(v)
		case domain.UserRole:
			role = v
		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		if role != domain.UserRoleAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
