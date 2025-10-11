package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"booknest/internal/domain"
)

const mockID = "cc0a35df-22f9-47cf-8b75-af0494343ede"

// mockUserService implements domain.UserService for tests
type mockUserService struct {
	users []domain.User
}

// Mock service functions
func (m *mockUserService) RegisterUser(ctx context.Context, in domain.UserInput) (string, error) {
	return "mock_token", nil
}

func (m *mockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// simulate success if id exists
	return nil
}

func (m *mockUserService) GetUsers(ctx context.Context) ([]domain.User, error) {
	return m.users, nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	for _, b := range m.users {
		if b.ID == id {
			return b, nil
		}
	}
	return domain.User{}, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, id uuid.UUID, entity domain.UserInput) (domain.User, error) {
	b := domain.User{
		ID:        id,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Password:  entity.Password,
		Email:     entity.Email,
	}
	return b, nil
}

func (m *mockUserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	return "mock_token", nil
}

// Tests
func TestRegisterUser_BadRequest(t *testing.T) {
	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestRegisterUser_Sucess(t *testing.T) {
	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	in := domain.UserInput{
		FirstName: "Test",
		LastName:  "Test",
		Password:  "Hashed",
		Email:     "test@email.com",
	}
	body, _ := json.Marshal(in)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.Token != "mock_token" {
		t.Fatalf("expected mock_token, got %s", resp.Token)
	}

	if resp.Message != "registration successful" {
		t.Fatalf("expected registration successful, got %s", resp.Message)
	}
}

func TestGetUsers_Empty(t *testing.T) {
	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var users []domain.User
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected zero users, got %d", len(users))
	}
}

func TestLogin_BadUUID(t *testing.T) {
	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad uuid, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestLogin_EmptyEmail(t *testing.T) {
	creds := make(map[string]string)
	creds["password"] = "hashed_password"
	body, _ := json.Marshal(creds)

	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin_EmptyPassword(t *testing.T) {
	creds := make(map[string]string)
	creds["email"] = "test1@email.com"
	body, _ := json.Marshal(creds)

	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	creds := make(map[string]string)
	creds["email"] = "test1@email.com"
	creds["password"] = "hashed_password"
	body, _ := json.Marshal(creds)

	svc := &mockUserService{users: []domain.User{}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.Token != "mock_token" {
		t.Fatalf("expected mock_token, got %s", resp.Token)
	}

	if resp.Message != "login successful" {
		t.Fatalf("expected login successful, got %s", resp.Message)
	}
}


func TestGetUsers_NonEmpty(t *testing.T) {
	// create service with one user
	b := domain.User{ID: uuid.MustParse(mockID),
		FirstName: "Test",
		LastName:  "Test",
		Password:  "Hashed",
		Email:     "test@email.com"}
	svc := &mockUserService{users: []domain.User{b}}
	r := setupRouter(nil, svc)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var users []domain.User
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected one user, got %d", len(users))
	}
}

func TestGetUser_BadUUID(t *testing.T) {
	r := setupRouter(nil, &mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/book/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad uuid, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetUser_Success(t *testing.T) {
	b := domain.User{
		ID:        uuid.MustParse(mockID),
		FirstName: "Test",
		LastName:  "Test",
		Password:  "Hashed",
		Email:     "test@email.com"}
	svc := &mockUserService{users: []domain.User{b}}

	r := setupRouter(nil, svc)
	req := httptest.NewRequest(http.MethodGet, "/user/cc0a35df-22f9-47cf-8b75-af0494343ede", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var user domain.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if user.ID.String() != mockID {
		t.Fatalf("expected id %s, got %s", mockID, user.ID.String())
	}
}

func TestDeleteUser_BadUUID(t *testing.T) {
	r := setupRouter(nil, &mockUserService{})
	req := httptest.NewRequest(http.MethodDelete, "/user/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad uuid, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteUser_Success(t *testing.T) {
	svc := &mockUserService{}
	r := setupRouter(nil, svc)
	req := httptest.NewRequest(http.MethodDelete, "/user/"+mockID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d, body: %s", w.Code, w.Body.String())
	}
}
