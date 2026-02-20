package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type authorRepo struct {
	gorm *gorm.DB
}

func NewAuthorRepo(gormDB *gorm.DB) domain.AuthorRepository {
	return &authorRepo{
		gorm: gormDB,
	}
}

func (r *authorRepo) FindByID(ctx context.Context, id uuid.UUID) (domain.Author, error) {
	var author domain.Author

	err := r.gorm.
		WithContext(ctx).
		Where("id = ?", id).
		First(&author).
		Error

	return author, err
}

func (r *authorRepo) FindByName(ctx context.Context, name string) (domain.Author, error) {
	var author domain.Author

	err := r.gorm.
		WithContext(ctx).
		Where("LOWER(name) = LOWER(?)", name).
		First(&author).
		Error

	return author, err
}

func (r *authorRepo) List(ctx context.Context, limit, offset int) ([]domain.Author, error) {
	var authors []domain.Author

	err := r.gorm.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&authors).Error

	return authors, err
}

func (r *authorRepo) Create(ctx context.Context, author *domain.Author) error {
	return r.gorm.WithContext(ctx).Create(author).Error
}

func (r *authorRepo) Update(ctx context.Context, author *domain.Author) error {
	return r.gorm.WithContext(ctx).Save(author).Error
}

func (r *authorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.gorm.WithContext(ctx).Delete(&domain.Author{}, "id = ?", id).Error
}
