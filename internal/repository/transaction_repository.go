package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func getTx(ctx context.Context) pgx.Tx {
    if ctx == nil {
        ctx = context.Background()
    }
    if txVal := ctx.Value(TxKey); txVal != nil {
        return txVal.(pgx.Tx)
    }
    return nil
}
