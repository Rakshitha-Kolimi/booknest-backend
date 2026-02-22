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

type mockOrderServiceController struct {
	checkoutFunc       func(ctx context.Context, userID uuid.UUID, input domain.CheckoutInput) (domain.OrderView, error)
	confirmPaymentFunc func(ctx context.Context, userID uuid.UUID, input domain.PaymentConfirmInput) (domain.OrderView, error)
	listUserOrdersFunc func(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.OrderView, error)
	listAllOrdersFunc  func(ctx context.Context, limit, offset int) ([]domain.OrderView, error)
}

func (m *mockOrderServiceController) Checkout(ctx context.Context, userID uuid.UUID, input domain.CheckoutInput) (domain.OrderView, error) {
	if m.checkoutFunc != nil {
		return m.checkoutFunc(ctx, userID, input)
	}
	return domain.OrderView{}, errors.New("not implemented")
}
func (m *mockOrderServiceController) ConfirmPayment(ctx context.Context, userID uuid.UUID, input domain.PaymentConfirmInput) (domain.OrderView, error) {
	if m.confirmPaymentFunc != nil {
		return m.confirmPaymentFunc(ctx, userID, input)
	}
	return domain.OrderView{}, errors.New("not implemented")
}
func (m *mockOrderServiceController) ListUserOrders(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.OrderView, error) {
	if m.listUserOrdersFunc != nil {
		return m.listUserOrdersFunc(ctx, userID, limit, offset)
	}
	return []domain.OrderView{}, nil
}
func (m *mockOrderServiceController) ListAllOrders(ctx context.Context, limit, offset int) ([]domain.OrderView, error) {
	if m.listAllOrdersFunc != nil {
		return m.listAllOrdersFunc(ctx, limit, offset)
	}
	return []domain.OrderView{}, nil
}

func TestOrderControllerCheckoutAndConfirm(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	orderID := uuid.New()
	calledCheckout := false
	calledConfirm := false

	svc := &mockOrderServiceController{
		checkoutFunc: func(ctx context.Context, gotUserID uuid.UUID, input domain.CheckoutInput) (domain.OrderView, error) {
			calledCheckout = true
			if gotUserID != userID || input.PaymentMethod != domain.PaymentCOD {
				t.Fatalf("unexpected checkout input")
			}
			return domain.OrderView{Order: domain.Order{ID: orderID}}, nil
		},
		confirmPaymentFunc: func(ctx context.Context, gotUserID uuid.UUID, input domain.PaymentConfirmInput) (domain.OrderView, error) {
			calledConfirm = true
			if gotUserID != userID || input.OrderID != orderID || !input.Success {
				t.Fatalf("unexpected confirm input")
			}
			return domain.OrderView{Order: domain.Order{ID: orderID}}, nil
		},
	}
	ctl := NewOrderController(svc).(*orderController)

	checkoutBody, _ := json.Marshal(domain.CheckoutInput{PaymentMethod: domain.PaymentCOD})
	cw := httptest.NewRecorder()
	cc, _ := gin.CreateTestContext(cw)
	cc.Set("user_id", userID.String())
	cc.Request = httptest.NewRequest(http.MethodPost, "/orders/checkout", bytes.NewBuffer(checkoutBody))
	cc.Request.Header.Set("Content-Type", "application/json")
	ctl.Checkout(cc)
	if cw.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", cw.Code)
	}

	confirmBody, _ := json.Marshal(domain.PaymentConfirmInput{OrderID: orderID, Success: true})
	rw := httptest.NewRecorder()
	rc, _ := gin.CreateTestContext(rw)
	rc.Set("user_id", userID.String())
	rc.Request = httptest.NewRequest(http.MethodPost, "/orders/confirm", bytes.NewBuffer(confirmBody))
	rc.Request.Header.Set("Content-Type", "application/json")
	ctl.ConfirmPayment(rc)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	if !calledCheckout || !calledConfirm {
		t.Fatalf("expected checkout and confirm handlers to call service")
	}
}

func TestOrderControllerListEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	svc := &mockOrderServiceController{
		listUserOrdersFunc: func(ctx context.Context, gotUserID uuid.UUID, limit, offset int) ([]domain.OrderView, error) {
			if gotUserID != userID || limit != 3 || offset != 1 {
				t.Fatalf("unexpected user list params")
			}
			return []domain.OrderView{}, nil
		},
		listAllOrdersFunc: func(ctx context.Context, limit, offset int) ([]domain.OrderView, error) {
			if limit != 4 || offset != 2 {
				t.Fatalf("unexpected admin list params")
			}
			return []domain.OrderView{}, nil
		},
	}
	ctl := NewOrderController(svc).(*orderController)

	uw := httptest.NewRecorder()
	uc, _ := gin.CreateTestContext(uw)
	uc.Set("user_id", userID.String())
	uc.Request = httptest.NewRequest(http.MethodGet, "/orders?limit=3&offset=1", nil)
	ctl.ListMyOrders(uc)
	if uw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", uw.Code)
	}

	aw := httptest.NewRecorder()
	ac, _ := gin.CreateTestContext(aw)
	ac.Request = httptest.NewRequest(http.MethodGet, "/admin/orders?limit=4&offset=2", nil)
	ctl.ListAllOrders(ac)
	if aw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", aw.Code)
	}
}

func TestOrderControllerUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctl := NewOrderController(&mockOrderServiceController{}).(*orderController)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/orders", nil)

	ctl.ListMyOrders(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
