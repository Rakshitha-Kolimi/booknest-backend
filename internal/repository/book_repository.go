package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

type bookRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewBookRepositoryImpl(db *pgxpool.Pool) domain.BookRepository {
	return &bookRepositoryImpl{db: db}
}

func (r *bookRepositoryImpl) GetBooks(ctx context.Context) ([]domain.Book, error) {
	tx := getTx(ctx)

	query := `
		SELECT id, title, author, price, stock, created_at, updated_at
		FROM books
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var rows pgx.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query(ctx, query)
	} else {
		rows, err = r.db.Query(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(
			&b.ID, &b.Title, &b.Author, &b.Price, &b.Stock,
			&b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

func (r *bookRepositoryImpl) GetBook(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	tx := getTx(ctx)

	query := `
		SELECT id, title, author, price, stock, created_at, updated_at
		FROM books
		WHERE id = $1 AND deleted_at IS NULL
	`

	var book domain.Book
	var err error

	if tx != nil {
		err = tx.QueryRow(ctx, query, id).Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Price,
			&book.Stock,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
	} else {
		err = r.db.QueryRow(ctx, query, id).Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Price,
			&book.Stock,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
	}

	return book, err
}

func (r *bookRepositoryImpl) CreateBook(ctx context.Context, entity *domain.Book) (err error) {
	// Check if the context has a transaction
	tx := getTx(ctx)

	q := `
		INSERT INTO books (title, author, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	args := []interface{}{entity.Title, entity.Author, entity.Price, entity.Stock}

	// Run the query
	if tx != nil {
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	}

	return err
}

func (r *bookRepositoryImpl) UpdateBook(ctx context.Context, entity *domain.Book) (err error) {
	tx := getTx(ctx)

	q := `
		UPDATE books SET title = $1, author=$2, price=$3, stock=$4,  updated_at=NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	args := []interface{}{entity.Title, entity.Author, entity.Price, entity.Stock, entity.ID}

	// Run the query
	if tx != nil {
		err = tx.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	} else {
		err = r.db.QueryRow(ctx, q, args...).Scan(&entity.UpdatedAt)
	}

	return err
}
func (r *bookRepositoryImpl) DeleteBook(ctx context.Context, id uuid.UUID) error {
	tx := getTx(ctx)

	query := `
		UPDATE books 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{id}

	var cmdTag pgconn.CommandTag
	var err error

	if tx != nil {
		cmdTag, err = tx.Exec(ctx, query, args...)
	} else {
		cmdTag, err = r.db.Exec(ctx, query, args...)
	}

	if err != nil {
		return err
	}

	// Check if a row was actually updated (book existed)
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("book not found or already deleted")
	}

	return nil
}
