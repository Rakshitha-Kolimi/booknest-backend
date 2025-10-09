package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"booknest/internal/domain"
)

type userServiceImpl struct {
	r domain.UserRepository
}

func NewUserServiceImpl(r domain.UserRepository) domain.UserService {
	return &userServiceImpl{r: r}
}

// RegisterUser handles new user registration and returns a JWT token
func (s *userServiceImpl) RegisterUser(ctx context.Context, in domain.UserInput) (string, error) {
	// check if user already exists (optional in your r)
	users, err := s.r.GetUsers(ctx)
	if err == nil {
		for _, u := range users {
			if u.Email == in.Email {
				return "", errors.New("email already registered")
			}
		}
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := domain.User{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		Password:  string(hashedPassword),
	}

	// save user
	if err := s.r.CreateUser(ctx, &user); err != nil {
		return "", err
	}

	// generate token
	token, err := generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// LoginUser validates credentials and returns a JWT token
func (s *userServiceImpl) LoginUser(ctx context.Context, email, password string) (string, error) {
	users, err := s.r.GetUsers(ctx)
	if err != nil {
		return "", err
	}

	var user *domain.User
	for _, u := range users {
		if u.Email == email {
			user = &u
			break
		}
	}

	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// generate token
	token, err := generateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUsers returns all users
func (s *userServiceImpl) GetUsers(ctx context.Context) ([]domain.User, error) {
	return s.r.GetUsers(ctx)
}

// GetUserByID returns a single user by ID
func (s *userServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.r.GetUser(ctx, id)
}

// DeleteUser performs a soft delete
func (s *userServiceImpl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.r.DeleteUser(ctx, id)
}

// Helper: Generate JWT token
func generateJWT(user domain.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "booknest_secret" // fallback for local dev
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
