package category_service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type categoryService struct {
	r domain.CategoryRepository
}

func NewCategoryService(r domain.CategoryRepository) domain.CategoryService {
	return &categoryService{
		r: r,
	}
}

func (s *categoryService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Category, error) {
	category, err := s.r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *categoryService) List(
	ctx context.Context,
	limit, offset int,
) ([]domain.Category, error) {
	return s.r.List(ctx, limit, offset)
}

func (s *categoryService) Create(
	ctx context.Context,
	input domain.CategoryInput,
) (*domain.Category, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("category name is required")
	}

	_, err := s.r.FindByName(ctx, name)
	if err == nil {
		return nil, errors.New("category name already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &domain.Category{
		ID:   uuid.New(),
		Name: name,
	}

	if err := s.r.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Update(
	ctx context.Context,
	id uuid.UUID,
	input domain.CategoryInput,
) (*domain.Category, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("category name is required")
	}

	category, err := s.r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing, err := s.r.FindByName(ctx, name)
	if err == nil && existing.ID != category.ID {
		return nil, errors.New("category name already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category.Name = name
	if err := s.r.Update(ctx, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}
