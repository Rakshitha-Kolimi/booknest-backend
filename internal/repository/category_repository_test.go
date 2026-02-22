package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

func TestCategoryRepo_CRUDAndList(t *testing.T) {
	db := setupTestDB(t, &domain.Category{})
	repo := &categoryRepo{gorm: db}

	ctx := context.Background()
	category := &domain.Category{ID: uuid.New(), Name: "Fiction"}
	require.NoError(t, repo.Create(ctx, category))

	foundByID, err := repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	require.Equal(t, category.Name, foundByID.Name)

	foundByName, err := repo.FindByName(ctx, "fiction")
	require.NoError(t, err)
	require.Equal(t, category.ID, foundByName.ID)

	second := &domain.Category{ID: uuid.New(), Name: "Science"}
	require.NoError(t, repo.Create(ctx, second))

	list, err := repo.List(ctx, 1, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "Fiction", list[0].Name)

	category.Name = "Fiction Updated"
	require.NoError(t, repo.Update(ctx, category))
	updated, err := repo.FindByID(ctx, category.ID)
	require.NoError(t, err)
	require.Equal(t, "Fiction Updated", updated.Name)

	require.NoError(t, repo.Delete(ctx, category.ID))
	_, err = repo.FindByID(ctx, category.ID)
	require.Error(t, err)
}
