package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"booknest/internal/domain"
)

func queryRowWithTx(
	ctx context.Context,
	pool domain.DBExecer,
	query string,
	args ...any,
) pgx.Row {
	// Check if there's an active transaction in the context
	// If so, use it to execute the query
	if txVal := ctx.Value(domain.TxKey); txVal != nil {
		if tx, ok := txVal.(pgx.Tx); ok {
			return tx.QueryRow(ctx, query, args...)
		}
	}

	// Otherwise, use the connection pool to execute the query
	return pool.QueryRow(ctx, query, args...)
}

func execWithTx(
	ctx context.Context,
	pool domain.DBExecer,
	query string,
	args ...any,
) error {
	// Check if there's an active transaction in the context
	// If so, use it to execute the query
	if txVal := ctx.Value(domain.TxKey); txVal != nil {
		if tx, ok := txVal.(pgx.Tx); ok {
			_, err := tx.Exec(ctx, query, args...)
			return err
		}
	}

	// Otherwise, use the connection pool to execute the query
	_, err := pool.Exec(ctx, query, args...)
	return err
}
