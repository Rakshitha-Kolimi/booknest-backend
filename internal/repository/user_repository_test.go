package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	return db
}

func TestUserRepo_FindByEmail(t *testing.T) {
	db := setupTestDB(t)

	user := domain.User{
		ID:        uuid.New(),
		Email:     "test@booknest.com",
		FirstName: "Test",
		IsActive:  true,
	}

	require.NoError(t, db.Create(&user).Error)

	repo := &userRepo{gorm: db}

	found, err := repo.FindByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.Equal(t, user.ID, found.ID)
	require.Equal(t, user.Email, found.Email)
}

func TestUserRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &userRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	user := &domain.User{
		ID:             uuid.New(),
		FirstName:      "Test",
		LastName:       "User",
		Email:          "test@booknest.com",
		Mobile:         "9999999999",
		Password:       "hashed",
		Role:           "user",
		IsActive:       true,
		EmailVerified:  false,
		MobileVerified: false,
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(user.ID, time.Now(), time.Now()),
		)

	err = repo.Create(context.Background(), user)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
