package util

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

var beginTx = func(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	return pool.Begin(ctx)
}

// Create a transaction function
// It should be common method to use in a database
func WithTransaction(
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(ctx context.Context) error,
) error {

	// Begin the transaction
	tx, err := beginTx(ctx, pool)
	if err != nil {
		return err
	}

	// Add value to context
	ctx = context.WithValue(ctx, domain.TxKey, tx)

	// If an error occurs, rollback
	if err := fn(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Commit the transaction
	return tx.Commit(ctx)
}
