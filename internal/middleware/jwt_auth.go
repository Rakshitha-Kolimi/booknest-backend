package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware verifies JWT token and injects user info into context
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract token from Authorization header
		authHeader := ctx.GetHeader("Authorization")

		// Expect header format: "Bearer <token>""
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// Return error if no Authorization header or invalid format
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			ctx.Abort()
			return
		}

		// Get the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "booknest_secret" // fallback for local dev
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		// Check for parsing errors or invalid token
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			ctx.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		// Attach user info to context for downstream use
		ctx.Set("user_id", claims["user_id"])
		ctx.Set("email", claims["email"])
		ctx.Set("user_role", claims["user_role"])

		ctx.Next()
	}
}
