package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"booknest/internal/domain"
)

// mockUserRepo is a lightweight in-memory mock implementing domain.UserRepository
type mockUserRepo struct {
	users       []domain.User
	getUsersErr error
	getUserRes  domain.User
	getUserErr  error
	createErr   error
	updateErr   error
	deleteErr   error

	// captures
	createdUser *domain.User
	deletedID   uuid.UUID
}

func (m *mockUserRepo) GetUsers(ctx context.Context) ([]domain.User, error) {
	if m.getUsersErr != nil {
		return nil, m.getUsersErr
	}
	return m.users, nil
}

func (m *mockUserRepo) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if m.getUserErr != nil {
		return domain.User{}, m.getUserErr
	}
	return m.getUserRes, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, entity *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	// simulate DB assigning ID and timestamps
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	m.createdUser = entity
	m.users = append(m.users, *entity)
	return nil
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, entity *domain.User) error {
	return m.updateErr
}

func (m *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	m.deletedID = id
	return m.deleteErr
}

func TestRegisterUser_SuccessAndDuplicate(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")

	// happy path
	repo := &mockUserRepo{users: []domain.User{}}
	svc := NewUserServiceImpl(repo)

	in := domain.UserInput{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     "alice@example.com",
		Password:  "password123",
	}

	tok, err := svc.RegisterUser(context.Background(), in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tok == "" {
		t.Fatalf("expected token, got empty string")
	}
	if repo.createdUser == nil {
		t.Fatalf("expected CreateUser to be called and createdUser captured")
	}
	// password should be hashed
	if bcrypt.CompareHashAndPassword([]byte(repo.createdUser.Password), []byte(in.Password)) != nil {
		t.Fatalf("stored password is not a valid bcrypt hash")
	}

	// duplicate email
	dupRepo := &mockUserRepo{users: []domain.User{{Email: "alice@example.com"}}}
	dupSvc := NewUserServiceImpl(dupRepo)
	if _, err := dupSvc.RegisterUser(context.Background(), in); err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}
}

func TestLoginUser_SuccessAndFailures(t *testing.T) {
	secretPassword := "mysecret"
	hashed, err := bcrypt.GenerateFromPassword([]byte(secretPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockUserRepo{users: []domain.User{{Email: "bob@example.com", Password: string(hashed)}}}
	svc := NewUserServiceImpl(repo)

	// success
	tok, err := svc.LoginUser(context.Background(), "bob@example.com", secretPassword)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tok == "" {
		t.Fatalf("expected token, got empty string")
	}

	// wrong password
	if _, err := svc.LoginUser(context.Background(), "bob@example.com", "wrongpwd"); err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}

	// unknown email
	if _, err := svc.LoginUser(context.Background(), "unknown@example.com", "whatever"); err == nil {
		t.Fatalf("expected error for unknown email, got nil")
	}
}

func TestGetUsers_GetUserByID_DeleteUser(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	u := domain.User{ID: id, FirstName: "C", LastName: "D", Email: "c@example.com", CreatedAt: now}

	repo := &mockUserRepo{users: []domain.User{u}, getUserRes: u}
	svc := NewUserServiceImpl(repo)

	list, err := svc.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("GetUsers error: %v", err)
	}
	if len(list) != 1 || list[0].Email != u.Email {
		t.Fatalf("unexpected users list: %+v", list)
	}

	got, err := svc.GetUserByID(context.Background(), id)
	if err != nil {
		t.Fatalf("GetUserByID error: %v", err)
	}
	if got.ID != id {
		t.Fatalf("expected id %s got %s", id, got.ID)
	}

	if err := svc.DeleteUser(context.Background(), id); err != nil {
		t.Fatalf("DeleteUser error: %v", err)
	}
	if repo.deletedID != id {
		t.Fatalf("expected deleted id %s, got %s", id, repo.deletedID)
	}
}
