package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"booknest/internal/domain"
)

type verificationTokenRepo struct {
	db   *sql.DB
	gorm *gorm.DB
	sb   squirrel.StatementBuilderType
}

func NewVerificationRepo(db *sql.DB, gormDB *gorm.DB) domain.VerificationTokenRepository {
	return &verificationTokenRepo{
		db:   db,
		gorm: gormDB,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *verificationTokenRepo) FindByUserIDAndType(
	ctx context.Context,
	userID uuid.UUID,
	tokenType domain.VerificationTokenType,
) (*domain.VerificationToken, error) {

	var token domain.VerificationToken

	err := r.gorm.
		WithContext(ctx).
		Where(
			"user_id = ? AND type = ? AND is_used = false",
			userID,
			tokenType,
		).
		Order("created_at DESC").
		First(&token).
		Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *verificationTokenRepo) Create(
	ctx context.Context,
	token *domain.VerificationToken,
) error {

	query, args, err := r.sb.
		Insert("verification_tokens").
		Columns(
			"user_id",
			"type",
			"token_hash",
			"expires_at",
			"is_used",
			"metadata",
		).
		Values(
			token.UserID,
			token.Type,
			token.TokenHash,
			token.ExpiresAt,
			token.IsUsed,
			token.Metadata,
		).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	row := queryRowWithTx(ctx, r.db, query, args...)

	return row.Scan(
		&token.ID,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
}

func (r *verificationTokenRepo) Update(
	ctx context.Context,
	token *domain.VerificationToken,
) error {

	query, args, err := r.sb.
		Update("verification_tokens").
		Set("is_used", token.IsUsed).
		Set("used_at", token.UsedAt).
		Set("expires_at", token.ExpiresAt).
		Set("metadata", token.Metadata).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": token.ID}).
		Suffix("RETURNING updated_at").
		ToSql()
	if err != nil {
		return err
	}

	row := queryRowWithTx(ctx, r.db, query, args...)
	return row.Scan(&token.UpdatedAt)
}

func (r *verificationTokenRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE verification_tokens
		SET deleted_at = NOW()
		WHERE id = $1
		RETURNING deleted_at
	`

	row := queryRowWithTx(ctx, r.db, query, id)

	var deletedAt time.Time
	return row.Scan(&deletedAt)
}
