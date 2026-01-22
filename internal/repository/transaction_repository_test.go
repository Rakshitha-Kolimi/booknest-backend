package repository

import (
	"context"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

func TestQueryRowWithTx_UsesTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.WithValue(context.Background(), domain.TxKey, mock)

	mock.ExpectQuery("SELECT 1").
		WillReturnRows(
			pgxmock.NewRows([]string{"col"}).AddRow(1),
		)

	row := queryRowWithTx(ctx, mock, "SELECT 1")

	var result int
	err = row.Scan(&result)

	require.NoError(t, err)
	require.Equal(t, 1, result)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestExecWithTx_UsesPool_WhenNoTx(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectExec("UPDATE users").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = execWithTx(context.Background(), mock, "UPDATE users")

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
