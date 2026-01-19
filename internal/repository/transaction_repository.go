package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

type TxKeyType string

const TxKey TxKeyType = "BookNest-Transactioner"

func getTx(ctx context.Context) pgx.Tx {
	if ctx == nil {
		ctx = context.Background()
	}
	if txVal := ctx.Value(TxKey); txVal != nil {
		return txVal.(pgx.Tx)
	}
	return nil
}

func queryRowWithTx(
	ctx context.Context,
	db *sql.DB,
	query string,
	args ...any,
) *sql.Row {

	if ctx == nil {
		ctx = context.Background()
	}

	if txVal := ctx.Value(TxKey); txVal != nil {
		return txVal.(*sql.Tx).QueryRowContext(ctx, query, args...)
	}

	return db.QueryRowContext(ctx, query, args...)
}
