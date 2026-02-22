package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type categoryRepo struct {
	gorm *gorm.DB
}

func NewCategoryRepo(gormDB *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{
		gorm: gormDB,
	}
}

func (r *categoryRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.Category, error) {
	var category domain.Category

	err := r.gorm.
		WithContext(ctx).
		Where("id = ?", id).
		First(&category).
		Error

	return category, err
}

func (r *categoryRepo) FindByName(ctx context.Context, name string) (domain.Category, error) {
	var category domain.Category

	err := r.gorm.
		WithContext(ctx).
		Where("LOWER(name) = LOWER(?)", name).
		First(&category).
		Error

	return category, err
}

func (r *categoryRepo) List(ctx context.Context, limit, offset int) ([]domain.Category, error) {
	var categories []domain.Category

	err := r.gorm.WithContext(ctx).
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *categoryRepo) Create(ctx context.Context, category *domain.Category) error {
	return r.gorm.WithContext(ctx).Create(category).Error
}

func (r *categoryRepo) Update(ctx context.Context, category *domain.Category) error {
	return r.gorm.WithContext(ctx).Save(category).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.gorm.WithContext(ctx).Delete(&domain.Category{}, "id = ?", id).Error
}
