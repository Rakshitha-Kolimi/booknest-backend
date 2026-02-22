package order_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"booknest/internal/domain"
)

type mockOrderRepository struct {
	listOrdersByUserFunc func(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.OrderView, error)
	listOrdersFunc       func(ctx context.Context, limit, offset int) ([]domain.OrderView, error)
}

func (m *mockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error { return nil }
func (m *mockOrderRepository) CreateOrderItems(ctx context.Context, items []domain.OrderItem) error {
	return nil
}
func (m *mockOrderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (domain.Order, error) {
	return domain.Order{}, errors.New("not implemented")
}
func (m *mockOrderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItemDetail, error) {
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepository) UpdateOrderPayment(ctx context.Context, orderID uuid.UUID, status domain.PaymentStatus, method domain.PaymentMethod) error {
	return nil
}
func (m *mockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return nil
}
func (m *mockOrderRepository) DecrementStock(ctx context.Context, items []domain.OrderItem) error {
	return nil
}
func (m *mockOrderRepository) ListOrdersByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.OrderView, error) {
	if m.listOrdersByUserFunc != nil {
		return m.listOrdersByUserFunc(ctx, userID, limit, offset)
	}
	return []domain.OrderView{}, nil
}
func (m *mockOrderRepository) ListOrders(ctx context.Context, limit, offset int) ([]domain.OrderView, error) {
	if m.listOrdersFunc != nil {
		return m.listOrdersFunc(ctx, limit, offset)
	}
	return []domain.OrderView{}, nil
}

type noopCartRepository struct{}

func (n *noopCartRepository) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (domain.Cart, error) {
	return domain.Cart{}, nil
}
func (n *noopCartRepository) GetCartItems(ctx context.Context, userID uuid.UUID) ([]domain.CartItemDetail, error) {
	return nil, nil
}
func (n *noopCartRepository) GetCartItemRecords(ctx context.Context, userID uuid.UUID) ([]domain.CartItemRecord, error) {
	return nil, nil
}
func (n *noopCartRepository) UpsertCartItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID, count int, unitPrice float64) error {
	return nil
}
func (n *noopCartRepository) RemoveCartItem(ctx context.Context, cartID uuid.UUID, bookID uuid.UUID) error {
	return nil
}
func (n *noopCartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error { return nil }

func TestPtrPaymentStatus(t *testing.T) {
	status := ptrPaymentStatus(domain.PaymentPaid)
	if status == nil || *status != domain.PaymentPaid {
		t.Fatalf("unexpected pointer status value: %+v", status)
	}
}

func TestOrderListPassThrough(t *testing.T) {
	userID := uuid.New()
	orders := []domain.OrderView{{Order: domain.Order{ID: uuid.New()}}}
	repo := &mockOrderRepository{
		listOrdersByUserFunc: func(ctx context.Context, gotUserID uuid.UUID, limit, offset int) ([]domain.OrderView, error) {
			if gotUserID != userID || limit != 10 || offset != 5 {
				t.Fatalf("unexpected user list params")
			}
			return orders, nil
		},
		listOrdersFunc: func(ctx context.Context, limit, offset int) ([]domain.OrderView, error) {
			if limit != 20 || offset != 10 {
				t.Fatalf("unexpected list params")
			}
			return orders, nil
		},
	}

	svc := NewOrderService(nil, repo, &noopCartRepository{})

	userOrders, err := svc.ListUserOrders(context.Background(), userID, 10, 5)
	if err != nil || len(userOrders) != 1 {
		t.Fatalf("unexpected ListUserOrders result: %+v, err=%v", userOrders, err)
	}

	allOrders, err := svc.ListAllOrders(context.Background(), 20, 10)
	if err != nil || len(allOrders) != 1 {
		t.Fatalf("unexpected ListAllOrders result: %+v, err=%v", allOrders, err)
	}
}
