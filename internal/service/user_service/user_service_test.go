package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"booknest/internal/domain"
)

// TestFindUser_Success tests successful user retrieval
func TestFindUser_Success(t *testing.T) {
	userID := uuid.New()
	expectedUser := domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Mobile:    "1234567890",
	}

	mockUserRepo := &MockUserRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.User, error) {
			if id == userID {
				return expectedUser, nil
			}
			return domain.User{}, errors.New("user not found")
		},
	}

	service := &userService{r: mockUserRepo}
	user, err := service.FindUser(context.Background(), userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, user.ID)
	}
}

// TestFindUser_NotFound tests user retrieval when user doesn't exist
func TestFindUser_NotFound(t *testing.T) {
	mockUserRepo := &MockUserRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		},
	}

	service := &userService{r: mockUserRepo}
	_, err := service.FindUser(context.Background(), uuid.New())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestLogin_SuccessByEmail tests successful login with email
func TestLogin_SuccessByEmail(t *testing.T) {
	userID := uuid.New()
	password := "password123"

	mockUserRepo := &MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			service := &userService{}
			hashedPassword := service.hashPassword(password)
			return domain.User{
				ID:       userID,
				Email:    email,
				Password: hashedPassword,
				Role:     domain.UserRoleUser,
			}, nil
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Email:    "test@example.com",
		Password: password,
	}

	token, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected non-empty token")
	}
}

// TestLogin_SuccessByMobile tests successful login with mobile
func TestLogin_SuccessByMobile(t *testing.T) {
	userID := uuid.New()
	password := "password123"

	mockUserRepo := &MockUserRepository{
		FindByMobileFunc: func(ctx context.Context, mobile string) (domain.User, error) {
			service := &userService{}
			hashedPassword := service.hashPassword(password)
			return domain.User{
				ID:       userID,
				Mobile:   mobile,
				Password: hashedPassword,
				Role:     domain.UserRoleUser,
			}, nil
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Mobile:   "1234567890",
		Password: password,
	}

	token, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected non-empty token")
	}
}

// TestLogin_InvalidPassword tests login with incorrect password
func TestLogin_InvalidPassword(t *testing.T) {
	password := "correctpassword123"

	mockUserRepo := &MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			service := &userService{}
			hashedPassword := service.hashPassword(password)
			return domain.User{
				ID:       uuid.New(),
				Email:    email,
				Password: hashedPassword,
				Role:     domain.UserRoleUser,
			}, nil
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := service.Login(context.Background(), input)

	if err == nil {
		t.Fatalf("expected error for invalid password")
	}

	if err.Error() != "invalid credentials" {
		t.Fatalf("expected 'invalid credentials' error, got %v", err)
	}
}

// TestLogin_UserNotFound tests login when user doesn't exist
func TestLogin_UserNotFound(t *testing.T) {
	mockUserRepo := &MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, err := service.Login(context.Background(), input)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestLogin_UpdatesLastLogin tests that login updates last_login timestamp
func TestLogin_UpdatesLastLogin(t *testing.T) {
	userID := uuid.New()
	password := "password123"
	updateCalled := false
	var capturedUser *domain.User

	mockUserRepo := &MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			service := &userService{}
			hashedPassword := service.hashPassword(password)
			return domain.User{
				ID:       userID,
				Email:    email,
				Password: hashedPassword,
				Role:     domain.UserRoleUser,
			}, nil
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			updateCalled = true
			capturedUser = user
			return nil
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Email:    "test@example.com",
		Password: password,
	}

	token, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	if !updateCalled {
		t.Fatalf("expected Update to be called")
	}

	if capturedUser.LastLogin == nil {
		t.Fatalf("expected LastLogin to be set")
	}
}

// TestResetPassword_Success tests successful password reset
// Note: This test requires a real DB pool and transaction support,
// so it's skipped and should be tested through integration tests
func TestResetPassword_Success(t *testing.T) {
	t.Skip("ResetPassword requires database transaction, tested through integration tests")
}

// TestDeleteUser_Success tests successful user deletion
func TestDeleteUser_Success(t *testing.T) {
	userID := uuid.New()
	deleteCalled := false

	mockUserRepo := &MockUserRepository{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			if id == userID {
				deleteCalled = true
				return nil
			}
			return errors.New("user not found")
		},
	}

	service := &userService{r: mockUserRepo}
	err := service.DeleteUser(context.Background(), userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !deleteCalled {
		t.Fatalf("expected Delete to be called")
	}
}

// TestDeleteUser_NotFound tests deletion of non-existent user
func TestDeleteUser_NotFound(t *testing.T) {
	mockUserRepo := &MockUserRepository{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return errors.New("user not found")
		},
	}

	service := &userService{r: mockUserRepo}
	err := service.DeleteUser(context.Background(), uuid.New())

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestVerifyEmail_CallsVerifyToken tests that VerifyEmail delegates to verifyToken
func TestVerifyEmail_CallsVerifyToken(t *testing.T) {
	// This test would require mocking the verifyToken method
	// Since verifyToken requires DB transaction, it's better tested through integration tests
	t.Skip("VerifyEmail requires database transaction, tested through integration tests")
}

// TestVerifyMobile_CallsVerifyToken tests that VerifyMobile delegates to verifyToken
func TestVerifyMobile_CallsVerifyToken(t *testing.T) {
	// This test would require mocking the verifyToken method
	// Since verifyToken requires DB transaction, it's better tested through integration tests
	t.Skip("VerifyMobile requires database transaction, tested through integration tests")
}

// TestResendEmailVerification_EmailAlreadyVerified tests resending when email already verified
// Note: This test requires a real DB pool and transaction support,
// so it's skipped and should be tested through integration tests
func TestResendEmailVerification_EmailAlreadyVerified(t *testing.T) {
	t.Skip("ResendEmailVerification requires database transaction, tested through integration tests")
}

// TestResendMobileOTP_MobileAlreadyVerified tests resending when mobile already verified
// Note: This test requires a real DB pool and transaction support,
// so it's skipped and should be tested through integration tests
func TestResendMobileOTP_MobileAlreadyVerified(t *testing.T) {
	t.Skip("ResendMobileOTP requires database transaction, tested through integration tests")
}

// TestNewUserService tests service initialization
func TestNewUserService(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockVerificationRepo := &MockVerificationTokenRepository{}

	service := NewUserService(nil, mockUserRepo, mockVerificationRepo)

	if service == nil {
		t.Fatalf("expected non-nil service")
	}

	userService, ok := service.(*userService)
	if !ok {
		t.Fatalf("expected *userService type")
	}

	if userService.r != mockUserRepo {
		t.Fatalf("expected user repository to be set")
	}

	if userService.vtr != mockVerificationRepo {
		t.Fatalf("expected verification token repository to be set")
	}
}

// TestLogin_TrimsContextualEmailAndMobile tests login prefers email when both provided
func TestLogin_PrefersEmailOverMobile(t *testing.T) {
	userID := uuid.New()
	password := "password123"
	emailCalled := false

	mockUserRepo := &MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (domain.User, error) {
			emailCalled = true
			service := &userService{}
			hashedPassword := service.hashPassword(password)
			return domain.User{
				ID:       userID,
				Email:    email,
				Password: hashedPassword,
				Role:     domain.UserRoleUser,
			}, nil
		},
		FindByMobileFunc: func(ctx context.Context, mobile string) (domain.User, error) {
			t.Fatalf("should not call FindByMobile when email is provided")
			return domain.User{}, errors.New("should not be called")
		},
		UpdateFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	service := &userService{r: mockUserRepo}
	input := domain.LoginInput{
		Email:    "test@example.com",
		Mobile:   "1234567890",
		Password: password,
	}

	_, err := service.Login(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !emailCalled {
		t.Fatalf("expected FindByEmail to be called")
	}
}

// TestFindUser_PassesContext tests that FindUser passes context properly
func TestFindUser_PassesContext(t *testing.T) {
	userID := uuid.New()
	contextPassed := false

	mockUserRepo := &MockUserRepository{
		FindByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.User, error) {
			if ctx == nil {
				t.Fatalf("expected non-nil context")
			}
			contextPassed = true
			return domain.User{ID: id}, nil
		},
	}

	service := &userService{r: mockUserRepo}
	user, err := service.FindUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, user.ID)
	}

	if !contextPassed {
		t.Fatalf("expected context to be passed to repository")
	}
}

// TestDeleteUser_PassesContext tests that DeleteUser passes context properly
func TestDeleteUser_PassesContext(t *testing.T) {
	userID := uuid.New()
	contextPassed := false

	mockUserRepo := &MockUserRepository{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			if ctx == nil {
				t.Fatalf("expected non-nil context")
			}
			contextPassed = true
			return nil
		},
	}

	service := &userService{r: mockUserRepo}
	err := service.DeleteUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !contextPassed {
		t.Fatalf("expected context to be passed to repository")
	}
}
