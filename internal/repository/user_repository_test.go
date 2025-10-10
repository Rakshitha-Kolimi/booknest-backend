package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

func TestGetUser(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "users")
	seed := seedUser(t)

	repo := NewUserRepositoryImpl(testDB)
	ctx := context.Background()

	got, err := repo.GetUser(ctx, seed.ID)
	require.NoError(t, err)
	require.Equal(t, seed.ID, got.ID)

	defer cleanTable(t, "users")
}

func TestGetUsers(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "users")
	seed := seedUser(t)

	repo := NewUserRepositoryImpl(testDB)
	ctx := context.Background()

	users, err := repo.GetUsers(ctx)
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, seed.ID, users[0].ID)
}

func TestCreateUser(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "users")
	repo := NewUserRepositoryImpl(testDB)
	ctx := context.Background()

	user := &domain.User{
		FirstName: "Test",
		LastName:  "Test",
		Password:  "Hashed",
		Email: "test@email.com",
	}

	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, user.ID)
}

func TestUpdateUser(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "users")
	seed := seedUser(t)

	repo := NewUserRepositoryImpl(testDB)
	ctx := context.Background()

	seed.LastName = "Test"
	err := repo.UpdateUser(ctx, &seed)
	require.NoError(t, err)

	updated, err := repo.GetUser(ctx, seed.ID)
	require.NoError(t, err)
	require.Equal(t, "Test", updated.LastName)
}

func TestDeleteUser(t *testing.T) {
	initTestDB(t)
	cleanTable(t, "users")
	seed := seedUser(t)

	repo := NewUserRepositoryImpl(testDB)
	ctx := context.Background()

	err := repo.DeleteUser(ctx, seed.ID)
	require.NoError(t, err)

	_, err = repo.GetUser(ctx, seed.ID)
	require.Error(t, err) // Should fail because user is soft-deleted
}
