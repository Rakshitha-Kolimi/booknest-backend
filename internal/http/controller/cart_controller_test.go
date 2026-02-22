package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"booknest/internal/domain"
)

type mockCartServiceController struct {
	getCartFunc    func(ctx context.Context, userID uuid.UUID) (domain.CartView, error)
	addItemFunc    func(ctx context.Context, userID uuid.UUID, input domain.CartItemInput) (domain.CartView, error)
	updateItemFunc func(ctx context.Context, userID uuid.UUID, input domain.CartItemInput) (domain.CartView, error)
	removeItemFunc func(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) (domain.CartView, error)
	clearFunc      func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockCartServiceController) GetCart(ctx context.Context, userID uuid.UUID) (domain.CartView, error) {
	if m.getCartFunc != nil {
		return m.getCartFunc(ctx, userID)
	}
	return domain.CartView{}, errors.New("not implemented")
}
func (m *mockCartServiceController) AddItem(ctx context.Context, userID uuid.UUID, input domain.CartItemInput) (domain.CartView, error) {
	if m.addItemFunc != nil {
		return m.addItemFunc(ctx, userID, input)
	}
	return domain.CartView{}, errors.New("not implemented")
}
func (m *mockCartServiceController) UpdateItem(ctx context.Context, userID uuid.UUID, input domain.CartItemInput) (domain.CartView, error) {
	if m.updateItemFunc != nil {
		return m.updateItemFunc(ctx, userID, input)
	}
	return domain.CartView{}, errors.New("not implemented")
}
func (m *mockCartServiceController) RemoveItem(ctx context.Context, userID uuid.UUID, bookID uuid.UUID) (domain.CartView, error) {
	if m.removeItemFunc != nil {
		return m.removeItemFunc(ctx, userID, bookID)
	}
	return domain.CartView{}, errors.New("not implemented")
}
func (m *mockCartServiceController) Clear(ctx context.Context, userID uuid.UUID) error {
	if m.clearFunc != nil {
		return m.clearFunc(ctx, userID)
	}
	return nil
}

func TestCartControllerGetCartUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctl := NewCartController(&mockCartServiceController{}).(*cartController)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cart", nil)

	ctl.GetCart(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCartControllerAddItemAndRemoveItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	bookID := uuid.New()
	svc := &mockCartServiceController{
		addItemFunc: func(ctx context.Context, gotUserID uuid.UUID, input domain.CartItemInput) (domain.CartView, error) {
			if gotUserID != userID || input.BookID != bookID || input.Count != 2 {
				t.Fatalf("unexpected add input")
			}
			return domain.CartView{TotalItems: 2}, nil
		},
		removeItemFunc: func(ctx context.Context, gotUserID uuid.UUID, gotBookID uuid.UUID) (domain.CartView, error) {
			if gotUserID != userID || gotBookID != bookID {
				t.Fatalf("unexpected remove input")
			}
			return domain.CartView{TotalItems: 0}, nil
		},
	}
	ctl := NewCartController(svc).(*cartController)

	body, _ := json.Marshal(domain.CartItemInput{BookID: bookID, Count: 2})
	aw := httptest.NewRecorder()
	ac, _ := gin.CreateTestContext(aw)
	ac.Set("user_id", userID.String())
	ac.Request = httptest.NewRequest(http.MethodPost, "/cart/items", bytes.NewBuffer(body))
	ac.Request.Header.Set("Content-Type", "application/json")
	ctl.AddItem(ac)
	if aw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", aw.Code)
	}

	rw := httptest.NewRecorder()
	rc, _ := gin.CreateTestContext(rw)
	rc.Set("user_id", userID.String())
	rc.Params = gin.Params{{Key: "book_id", Value: bookID.String()}}
	rc.Request = httptest.NewRequest(http.MethodDelete, "/cart/items/"+bookID.String(), nil)
	ctl.RemoveItem(rc)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestCartControllerClearCart(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	called := false
	svc := &mockCartServiceController{clearFunc: func(ctx context.Context, gotUserID uuid.UUID) error {
		called = true
		if gotUserID != userID {
			t.Fatalf("unexpected user id")
		}
		return nil
	}}
	ctl := NewCartController(svc).(*cartController)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", userID.String())
	c.Request = httptest.NewRequest(http.MethodPost, "/cart/clear", nil)

	ctl.ClearCart(c)
	if !called {
		t.Fatalf("expected clear to be called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
