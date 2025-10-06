package repository

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

type TxKeyType string

const TxKey TxKeyType = "BookNest-Transactioner"

type bookRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewBookRepositoryImpl(db *pgxpool.Pool) domain.BookRepository {
	return &bookRepositoryImpl{db: db}
}

func (r *bookRepositoryImpl) CreateBook(ctx context.Context, entity *domain.Book) (err error){
	// Check if the context has a transaction
	if ctx == nil {
		ctx = context.Background()
	}
	txVal := ctx.Value(TxKey)

	q := `
		INSERT INTO books (title, author, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	args := []interface{}{entity.Title, entity.Author, entity.Price, entity.Stock}

	// Run the query
	if txVal != nil {
		tx := txVal.(pgx.Tx)
		err = tx.QueryRow(q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	}

	return err
}
