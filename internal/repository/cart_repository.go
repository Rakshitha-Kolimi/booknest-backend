package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"booknest/internal/domain"
)

type cartRepo struct {
	db domain.DBExecer
	sb squirrel.StatementBuilderType
}

func NewCartRepo(db *pgxpool.Pool) domain.CartRepository {
	return &cartRepo{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *cartRepo) GetOrCreateCart(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Cart, error) {
	var cart domain.Cart

	query := `
		INSERT INTO carts (id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
		RETURNING id, user_id;
	`

	row := queryRowWithTx(ctx, r.db, query, uuid.New(), userID)
	err := row.Scan(&cart.ID, &cart.UserID)
	return cart, err
}

func (r *cartRepo) GetCartItems(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CartItemDetail, error) {
	query := `
		SELECT
			ci.book_id,
			b.name,
			b.author_name,
			b.image_url,
			(b.price - (b.price * b.discount_percentage / 100)) AS unit_price,
			ci.count,
			((b.price - (b.price * b.discount_percentage / 100)) * ci.count) AS line_total
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id AND ci.deleted_at IS NULL
		JOIN books b ON b.id = ci.book_id AND b.deleted_at IS NULL
		WHERE c.user_id = $1
		ORDER BY ci.created_at DESC;
	`

	rows, err := queryWithTx(ctx, r.db, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.CartItemDetail, 0)
	for rows.Next() {
		var item domain.CartItemDetail
		if err := rows.Scan(
			&item.BookID,
			&item.Name,
			&item.AuthorName,
			&item.ImageURL,
			&item.UnitPrice,
			&item.Count,
			&item.LineTotal,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *cartRepo) GetCartItemRecords(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.CartItemRecord, error) {
	query := `
		SELECT
			ci.book_id,
			ci.count,
			(b.price - (b.price * b.discount_percentage / 100)) AS unit_price,
			b.available_stock
		FROM carts c
		JOIN cart_items ci ON ci.cart_id = c.id AND ci.deleted_at IS NULL
		JOIN books b ON b.id = ci.book_id AND b.deleted_at IS NULL
		WHERE c.user_id = $1
		ORDER BY ci.created_at DESC;
	`

	rows, err := queryWithTx(ctx, r.db, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.CartItemRecord, 0)
	for rows.Next() {
		var item domain.CartItemRecord
		if err := rows.Scan(
			&item.BookID,
			&item.Count,
			&item.UnitPrice,
			&item.AvailableStock,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *cartRepo) UpsertCartItem(
	ctx context.Context,
	cartID uuid.UUID,
	bookID uuid.UUID,
	count int,
	unitPrice float64,
) error {
	query := `
		INSERT INTO cart_items (cart_id, book_id, count, cart_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (cart_id, book_id) DO UPDATE
		SET count = EXCLUDED.count,
		    cart_price = EXCLUDED.cart_price,
		    updated_at = NOW(),
		    deleted_at = NULL;
	`

	return execWithTx(ctx, r.db, query, cartID, bookID, count, unitPrice)
}

func (r *cartRepo) RemoveCartItem(
	ctx context.Context,
	cartID uuid.UUID,
	bookID uuid.UUID,
) error {
	query := `
		UPDATE cart_items
		SET deleted_at = NOW()
		WHERE cart_id = $1 AND book_id = $2 AND deleted_at IS NULL;
	`
	return execWithTx(ctx, r.db, query, cartID, bookID)
}

func (r *cartRepo) ClearCart(
	ctx context.Context,
	cartID uuid.UUID,
) error {
	query := `
		UPDATE cart_items
		SET deleted_at = NOW()
		WHERE cart_id = $1 AND deleted_at IS NULL;
	`
	return execWithTx(ctx, r.db, query, cartID)
}
