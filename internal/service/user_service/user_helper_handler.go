package user_service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"booknest/internal/domain"
)

func (s userService) hashPassword(p string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		slog.Error("Cannot hash the password", "error", err)
		return ""
	}

	return string(hashed)
}

func (s userService) comparePassword(hash, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}

func (s userService) generateRawToken() (string, error) {
	b := make([]byte, 32) // 256-bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s userService) generateOTP(length int) string {
	const digits = "1234567890"
	if length <= 0 {
		return ""
	}
	buffer := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, buffer, length)
	if n != length || err != nil {
		return ""
	}
	for i := 0; i < len(buffer); i++ {
		buffer[i] = digits[int(buffer[i])%len(digits)]
	}
	return string(buffer)
}

func (s userService) generateTokenHash(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

func (s userService) generateJWT(user domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "booknest_secret" // fallback for local dev
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"user_role": user.Role,
		"email":     user.Email,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
