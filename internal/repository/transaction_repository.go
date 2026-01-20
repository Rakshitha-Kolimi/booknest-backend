package repository

import (
	"context"

	"booknest/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func queryRowWithTx(
	ctx context.Context,
	pool *pgxpool.Pool,
	query string,
	args ...any,
) pgx.Row {

	if txVal := ctx.Value(domain.TxKey); txVal != nil {
		if tx, ok := txVal.(pgx.Tx); ok {
			return tx.QueryRow(ctx, query, args...)
		}
	}

	return pool.QueryRow(ctx, query, args...)
}
