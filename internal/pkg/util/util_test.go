package util

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxmock "github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"booknest/internal/domain"
)

func TestWithTransaction_CommitsOnSuccess(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	originalBegin := beginTx
	defer func() { beginTx = originalBegin }()

	mockPool.ExpectBegin()
	mockPool.ExpectCommit()

	beginTx = func(ctx context.Context, _ *pgxpool.Pool) (pgx.Tx, error) {
		return mockPool.Begin(ctx)
	}

	var txInContext any
	err = WithTransaction(context.Background(), nil, func(txCtx context.Context) error {
		txInContext = txCtx.Value(domain.TxKey)
		return nil
	})

	require.NoError(t, err)
	require.NotNil(t, txInContext)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestWithTransaction_RollsBackOnCallbackError(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPool.Close()

	originalBegin := beginTx
	defer func() { beginTx = originalBegin }()

	mockPool.ExpectBegin()
	mockPool.ExpectRollback()

	beginTx = func(ctx context.Context, _ *pgxpool.Pool) (pgx.Tx, error) {
		return mockPool.Begin(ctx)
	}

	expectedErr := errors.New("callback failed")
	err = WithTransaction(context.Background(), nil, func(txCtx context.Context) error {
		return expectedErr
	})

	require.ErrorIs(t, err, expectedErr)
	require.NoError(t, mockPool.ExpectationsWereMet())
}

func TestWithTransaction_BeginError(t *testing.T) {
	originalBegin := beginTx
	defer func() { beginTx = originalBegin }()

	expectedErr := errors.New("begin failed")
	beginTx = func(ctx context.Context, _ *pgxpool.Pool) (pgx.Tx, error) {
		return nil, expectedErr
	}

	called := false
	err := WithTransaction(context.Background(), nil, func(txCtx context.Context) error {
		called = true
		return nil
	})

	require.ErrorIs(t, err, expectedErr)
	require.False(t, called)
}
