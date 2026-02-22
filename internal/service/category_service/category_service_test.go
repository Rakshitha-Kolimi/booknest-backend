package category_service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type mockCategoryRepository struct {
	findByIDFunc   func(ctx context.Context, id uuid.UUID) (domain.Category, error)
	findByNameFunc func(ctx context.Context, name string) (domain.Category, error)
	listFunc       func(ctx context.Context, limit, offset int) ([]domain.Category, error)
	createFunc     func(ctx context.Context, category *domain.Category) error
	updateFunc     func(ctx context.Context, category *domain.Category) error
	deleteFunc     func(ctx context.Context, id uuid.UUID) error
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return domain.Category{}, errors.New("not implemented")
}

func (m *mockCategoryRepository) FindByName(ctx context.Context, name string) (domain.Category, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(ctx, name)
	}
	return domain.Category{}, gorm.ErrRecordNotFound
}

func (m *mockCategoryRepository) List(ctx context.Context, limit, offset int) ([]domain.Category, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return []domain.Category{}, nil
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, category)
	}
	return nil
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, category)
	}
	return nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestCreateCategorySuccess(t *testing.T) {
	repo := &mockCategoryRepository{
		findByNameFunc: func(ctx context.Context, name string) (domain.Category, error) {
			if name != "Fiction" {
				t.Fatalf("expected trimmed name, got %q", name)
			}
			return domain.Category{}, gorm.ErrRecordNotFound
		},
		createFunc: func(ctx context.Context, category *domain.Category) error {
			if category.ID == uuid.Nil {
				t.Fatalf("expected non-nil id")
			}
			if category.Name != "Fiction" {
				t.Fatalf("unexpected category name: %q", category.Name)
			}
			return nil
		},
	}

	svc := NewCategoryService(repo)
	category, err := svc.Create(context.Background(), domain.CategoryInput{Name: " Fiction "})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if category == nil || category.Name != "Fiction" {
		t.Fatalf("unexpected category: %+v", category)
	}
}

func TestCreateCategoryValidationAndDuplicate(t *testing.T) {
	svc := NewCategoryService(&mockCategoryRepository{})
	_, err := svc.Create(context.Background(), domain.CategoryInput{Name: " "})
	if err == nil || err.Error() != "category name is required" {
		t.Fatalf("expected required-name error, got %v", err)
	}

	repo := &mockCategoryRepository{
		findByNameFunc: func(ctx context.Context, name string) (domain.Category, error) {
			return domain.Category{ID: uuid.New(), Name: name}, nil
		},
	}
	svc = NewCategoryService(repo)
	_, err = svc.Create(context.Background(), domain.CategoryInput{Name: "Fiction"})
	if err == nil || err.Error() != "category name already exists" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestUpdateCategorySuccessAndConflict(t *testing.T) {
	categoryID := uuid.New()
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Category, error) {
			return domain.Category{ID: categoryID, Name: "Old"}, nil
		},
		findByNameFunc: func(ctx context.Context, name string) (domain.Category, error) {
			if name != "New Name" {
				t.Fatalf("expected trimmed update name, got %q", name)
			}
			return domain.Category{}, gorm.ErrRecordNotFound
		},
		updateFunc: func(ctx context.Context, category *domain.Category) error {
			if category.ID != categoryID || category.Name != "New Name" {
				t.Fatalf("unexpected update payload: %+v", category)
			}
			return nil
		},
	}

	svc := NewCategoryService(repo)
	updated, err := svc.Update(context.Background(), categoryID, domain.CategoryInput{Name: " New Name "})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "New Name" {
		t.Fatalf("unexpected updated category: %+v", updated)
	}

	conflictRepo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Category, error) {
			return domain.Category{ID: categoryID, Name: "Old"}, nil
		},
		findByNameFunc: func(ctx context.Context, name string) (domain.Category, error) {
			return domain.Category{ID: uuid.New(), Name: name}, nil
		},
	}

	svc = NewCategoryService(conflictRepo)
	_, err = svc.Update(context.Background(), categoryID, domain.CategoryInput{Name: "New Name"})
	if err == nil || err.Error() != "category name already exists" {
		t.Fatalf("expected duplicate-name error, got %v", err)
	}
}

func TestCategoryReadAndDeletePassThrough(t *testing.T) {
	categoryID := uuid.New()
	expected := domain.Category{ID: categoryID, Name: "Fiction"}
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Category, error) {
			return expected, nil
		},
		listFunc: func(ctx context.Context, limit, offset int) ([]domain.Category, error) {
			if limit != 5 || offset != 10 {
				t.Fatalf("unexpected pagination: %d/%d", limit, offset)
			}
			return []domain.Category{expected}, nil
		},
		deleteFunc: func(ctx context.Context, id uuid.UUID) error {
			if id != categoryID {
				t.Fatalf("unexpected id: %s", id)
			}
			return nil
		},
	}

	svc := NewCategoryService(repo)
	category, err := svc.FindByID(context.Background(), categoryID)
	if err != nil || category.ID != categoryID {
		t.Fatalf("unexpected find result: %+v, err=%v", category, err)
	}

	list, err := svc.List(context.Background(), 5, 10)
	if err != nil || len(list) != 1 {
		t.Fatalf("unexpected list result: %+v, err=%v", list, err)
	}

	if err := svc.Delete(context.Background(), categoryID); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
}
