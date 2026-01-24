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

func TestVerificationRepo_FindByHashAndType(t *testing.T) {
	db := setupTestDB(t, &domain.VerificationToken{})

	vt := domain.VerificationToken{
		ID:        uuid.New(),
		TokenHash: "token_hash",
		Type:      domain.VerificationEmail,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	require.NoError(t, db.Create(&vt).Error)

	repo := &verificationTokenRepo{gorm: db}

	found, err := repo.FindByHashAndType(context.Background(), vt.TokenHash, vt.Type)

	require.NoError(t, err)
	require.Equal(t, vt.ID, found.ID)
	require.Equal(t, vt.TokenHash, found.TokenHash)
	require.Equal(t, vt.Type, found.Type)
}

func TestVerificationRepo_FindByUserIDAndType(t *testing.T) {
	db := setupTestDB(t, &domain.VerificationToken{})

	vt := domain.VerificationToken{
		ID:        uuid.New(),
		TokenHash: "token_hash",
		Type:      domain.VerificationEmail,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		UserID:    uuid.New(),
	}

	require.NoError(t, db.Create(&vt).Error)

	repo := &verificationTokenRepo{gorm: db}

	found, err := repo.FindByUserIDAndType(context.TODO(), vt.UserID, vt.Type)

	require.NoError(t, err)
	require.Equal(t, vt.ID, found.ID)
	require.Equal(t, vt.TokenHash, found.TokenHash)
	require.Equal(t, vt.Type, found.Type)
	require.Equal(t, vt.UserID, found.UserID)
}

func TestVerificationRepo_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &verificationTokenRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	vt := &domain.VerificationToken{
		ID:        uuid.New(),
		TokenHash: "token_hash",
		Type:      domain.VerificationEmail,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		UserID:    uuid.New(),
	}

	mock.ExpectQuery("INSERT INTO verification_tokens").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(), pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(vt.ID, time.Now(), time.Now()),
		)

	err = repo.Create(context.Background(), vt)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepo_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &verificationTokenRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	vt := &domain.VerificationToken{
		ID:        uuid.New(),
		TokenHash: "token_hash",
		Type:      domain.VerificationEmail,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		UserID:    uuid.New(),
	}

	mock.ExpectQuery("UPDATE verification_tokens").
		WithArgs(
			pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
			pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"updated_at"}).
				AddRow(time.Now()),
		)

	err = repo.Update(context.Background(), vt)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepo_InvalidateByUserAndType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &verificationTokenRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	vt := &domain.VerificationToken{
		UserID: uuid.New(),
	}

	mock.ExpectExec(`UPDATE verification_tokens`).
		WithArgs(
			true, // is_used
			false,
			domain.VerificationEmail,
			vt.UserID.String(),
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.InvalidateByUserAndType(context.Background(), vt.UserID, domain.VerificationEmail)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepo_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &verificationTokenRepo{
		db: mock,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	 vt := &domain.VerificationToken{
		ID: uuid.New(),
	}

	mock.ExpectQuery("UPDATE verification_tokens").
		WithArgs(
			pgxmock.AnyArg(),
		).
		WillReturnRows(
			pgxmock.NewRows([]string{"deleted_at"}).
				AddRow(time.Now()),
		)

	err = repo.Delete(context.Background(), vt.ID)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
