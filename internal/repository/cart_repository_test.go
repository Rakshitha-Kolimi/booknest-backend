package repository

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestCartRepo_GetOrCreateCart(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &cartRepo{db: mock, sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	userID := uuid.New()
	cartID := uuid.New()

	mock.ExpectQuery("WITH inserted AS").
		WithArgs(pgxmock.AnyArg(), userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id"}).AddRow(cartID, userID))

	cart, err := repo.GetOrCreateCart(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, cartID, cart.ID)
	require.Equal(t, userID, cart.UserID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCartRepo_UpsertAndRemoveAndClear(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &cartRepo{db: mock, sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
	cartID := uuid.New()
	bookID := uuid.New()

	mock.ExpectExec("INSERT INTO cart_items").
		WithArgs(cartID, bookID, 2, 99.5).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	require.NoError(t, repo.UpsertCartItem(context.Background(), cartID, bookID, 2, 99.5))

	mock.ExpectExec("UPDATE cart_items").
		WithArgs(cartID, bookID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	require.NoError(t, repo.RemoveCartItem(context.Background(), cartID, bookID))

	mock.ExpectExec("UPDATE cart_items").
		WithArgs(cartID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 2))
	require.NoError(t, repo.ClearCart(context.Background(), cartID))

	require.NoError(t, mock.ExpectationsWereMet())
}
