package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

// mockBookService implements domain.BookService for tests
type mockBookService struct {
	books []domain.Book
}

// Mock service functions
func (m *mockBookService) CreateBook(ctx context.Context, in domain.BookInput) (domain.Book, error) {
	b := domain.Book{
		ID:     uuid.New(),
		Title:  in.Title,
		Author: in.Author,
		Price:  in.Price,
		Stock:  in.Stock,
	}
	m.books = append(m.books, b)
	return b, nil
}

func (m *mockBookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	// simulate success if id exists
	return nil
}

func (m *mockBookService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.books, nil
}

func (m *mockBookService) GetBook(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	for _, b := range m.books {
		if b.ID == id {
			return b, nil
		}
	}
	return domain.Book{}, nil
}

func (m *mockBookService) UpdateBook(ctx context.Context, id uuid.UUID, entity domain.BookInput) (domain.Book, error) {
	b := domain.Book{
		ID:     id,
		Title:  entity.Title,
		Author: entity.Author,
		Price:  entity.Price,
		Stock:  entity.Stock,
	}
	return b, nil
}

func setupRouter(bs domain.BookService, us domain.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	bc := NewBookController(bs)
	uc := NewUserController(us)

	r.GET("/health", GetHealth)
	r.GET("/books", bc.GetBooks)
	r.GET("/book/:id", bc.GetBook)
	r.PUT("/book/:id", bc.UpdateBook)
	r.POST("/books", bc.AddBook)
	r.DELETE("/book/:id", bc.DeleteBook)

	r.GET("/users", uc.GetUsers)
	r.GET("/user/:id", uc.GetUserByID)
	r.POST("/login", uc.LoginUser)
	r.POST("/register", uc.RegisterUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	return r
}

func TestGetHealth(t *testing.T) {
	r := setupRouter(&mockBookService{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("unexpected body: %v", resp)
	}
}

func TestGetBooks_Empty(t *testing.T) {
	svc := &mockBookService{books: []domain.Book{}}
	r := setupRouter(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &books); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("expected zero books, got %d", len(books))
	}
}

func TestAddBookAndGetBook(t *testing.T) {
	svc := &mockBookService{books: []domain.Book{}}
	r := setupRouter(svc, nil)

	in := domain.BookInput{Title: "Go in Action", Author: "William", Price: 39.99, Stock: 10}
	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on create, got %d, body: %s", w.Code, w.Body.String())
	}
	var created domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if created.Title != in.Title {
		t.Fatalf("title mismatch: expected %s, got %s", in.Title, created.Title)
	}

	// now fetch by id
	url := "/book/" + created.ID.String()
	req = httptest.NewRequest(http.MethodGet, url, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on get, got %d, body: %s", w.Code, w.Body.String())
	}
	var fetched domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if fetched.ID != created.ID {
		t.Fatalf("id mismatch: expected %s, got %s", created.ID, fetched.ID)
	}
}

func TestGetBook_BadUUID(t *testing.T) {
	r := setupRouter(&mockBookService{}, nil)
	req := httptest.NewRequest(http.MethodGet, "/book/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad uuid, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteBook_Success(t *testing.T) {
	svc := &mockBookService{}
	r := setupRouter(svc, nil)
	id := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/book/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetBooks_NonEmpty(t *testing.T) {
	// create service with one book
	b := domain.Book{ID: uuid.New(), Title: "Test", Author: "Author", Price: 9.99, Stock: 1}
	svc := &mockBookService{books: []domain.Book{b}}
	r := setupRouter(svc, nil)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var books []domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &books); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("expected one book, got %d", len(books))
	}
}

func TestUpdateBook_Success(t *testing.T) {
	svc := &mockBookService{books: []domain.Book{}}
	r := setupRouter(svc, nil)

	// create a book first
	in := domain.BookInput{Title: "Original", Author: "A", Price: 10.0, Stock: 5}
	body, _ := json.Marshal(in)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create failed: %d %s", w.Code, w.Body.String())
	}
	var created domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	// update
	upd := domain.BookInput{Title: "Updated", Author: "B", Price: 12.5, Stock: 3}
	ub, _ := json.Marshal(upd)
	req = httptest.NewRequest(http.MethodPut, "/book/"+created.ID.String(), bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d, body: %s", w.Code, w.Body.String())
	}
	var updated domain.Book
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if updated.Title != upd.Title {
		t.Fatalf("expected title %s, got %s", upd.Title, updated.Title)
	}
}

func TestUpdateBook_BadUUID(t *testing.T) {
	r := setupRouter(&mockBookService{}, nil)
	upd := domain.BookInput{Title: "X"}
	body, _ := json.Marshal(upd)
	req := httptest.NewRequest(http.MethodPut, "/book/not-a-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad uuid, got %d, body: %s", w.Code, w.Body.String())
	}
}
