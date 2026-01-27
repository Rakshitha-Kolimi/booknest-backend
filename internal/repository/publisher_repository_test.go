package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

func TestPublisherRepo_FindByID(t *testing.T) {
	db := setupTestDB(t, &domain.Publisher{})

	publisher := domain.Publisher{
		ID:          uuid.New(),
	}

	require.NoError(t, db.Create(&publisher).Error)

	repo := &publisherRepo{gorm: db}

	found, err := repo.FindByID(context.TODO(), publisher.ID)

	require.NoError(t, err)
	require.Equal(t, publisher.ID, found.ID)
}

func TestPublisherRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &publisherRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	publisher := &domain.Publisher{
		ID:          uuid.New(),
		LegalName:   "Legal",
		TradingName: "Trading",
		Email:       "test@mail.com",
		Mobile:      "+911234567890",
		Address:     "Addr",
		City:        "City",
		State:       "State",
		Country:     "Country",
		Zipcode:     "123456",
		IsActive:    true,
	}

	mock.ExpectQuery("INSERT INTO publishers").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(publisher.ID, time.Now(), time.Now()),
		)

	err = repo.Create(context.Background(), publisher)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPublisherRepo_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &publisherRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	publisher := &domain.Publisher{
		ID:          uuid.New(),
		LegalName:   "Legal",
		TradingName: "Trading",
		Email:       "test@mail.com",
		Mobile:      "+911234567890",
		Address:     "Addr",
		City:        "City",
		State:       "State",
		Country:     "Country",
		Zipcode:     "123456",
		IsActive:    true,
	}

	mock.ExpectQuery("UPDATE publishers").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"updated_at"}).
				AddRow(time.Now()),
		)

	err = repo.Update(context.Background(), publisher)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPublisherRepo_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &publisherRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	publisher := &domain.Publisher{
		ID: uuid.New(),
	}

	mock.ExpectQuery("UPDATE publishers").
		WithArgs(
			pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"deleted_at"}).
				AddRow(time.Now()),
		)

	err = repo.Delete(context.Background(), publisher.ID)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
