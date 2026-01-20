package user_service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
	"booknest/internal/pkg/util"
)

type userService struct {
	db  *pgxpool.Pool
	r   domain.UserRepository
	vtr domain.VerificationTokenRepository
}

func NewUserService(db *pgxpool.Pool, r domain.UserRepository, vtr domain.VerificationTokenRepository) domain.UserService {
	return &userService{
		db:  db,
		r:   r,
		vtr: vtr,
	}
}

func (s *userService) FindUser(
	ctx context.Context,
	id uuid.UUID,
) (domain.User, error) {
	return s.r.FindByID(ctx, id)
}

func (s *userService) Register(
	ctx context.Context,
	in domain.UserInput,
) error {

	return util.WithTransaction(ctx, s.db, func(txCtx context.Context) error {
		// Create an user domain
		user := &domain.User{
			ID:        uuid.New(),
			FirstName: in.FirstName,
			LastName:  in.LastName,
			Email:     in.Email,
			Mobile:    in.Mobile,
			Password:  s.hashPassword(in.Password),
			Role:      in.Role,
			IsActive:  true,
		}

		// Create the user
		if err := s.r.Create(txCtx, user); err != nil {
			return err
		}

		// Create Verification code to send the email
		verificationCode, err := s.generateRawToken()
		if err != nil {
			return err
		}
		token := &domain.VerificationToken{
			UserID:    user.ID,
			Type:      domain.VerificationEmail,
			TokenHash: s.generateTokenHash(verificationCode),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		// Add token to the database
		return s.vtr.Create(txCtx, token)
	})
}

func (s *userService) Login(
	ctx context.Context,
	in domain.LoginInput,
) (string, error) {

	var user domain.User
	var err error

	// Get user by email or mobile
	if in.Email != "" {
		user, err = s.r.FindByEmail(ctx, in.Email)
	} else {
		user, err = s.r.FindByMobile(ctx, in.Mobile)
	}
	if err != nil {
		return "", err
	}

	// Validate the password
	if !s.comparePassword(user.Password, in.Password) {
		return "", errors.New("invalid credentials")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now

	// Update user
	if err := s.r.Update(ctx, &user); err != nil {
		return "", err
	}

	// Generate JWT token
	return s.generateJWT(user)
}

func (s *userService) ResetPassword(
	ctx context.Context,
	userID uuid.UUID,
	newPassword string,
) error {

	return util.WithTransaction(ctx, s.db, func(txCtx context.Context) error {
		// Get user by ID
		user, err := s.r.FindByID(txCtx, userID)
		if err != nil {
			return err
		}

		// Add hashed password to the user
		user.Password = s.hashPassword(newPassword)

		// Update the user
		return s.r.Update(txCtx, &user)
	})
}

func (s *userService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.r.Delete(ctx, id)
}
