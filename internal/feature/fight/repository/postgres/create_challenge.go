package fight_pg_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

const uniqueViolationCode = "23505"
const checkViolationCode = "23514"

// CreateChallenge списывает ставку у создателя и создаёт вызов одной транзакцией:
// мёд «в банке» вызова с момента создания.
func (r *Repository) CreateChallenge(
	ctx context.Context,
	creatorUserID int,
	amount int64,
) (domain.ActiveChallenge, error) {
	const lockHoneyQuery = `
		SELECT honey FROM honey.user_honey WHERE user_id = $1 FOR UPDATE
	`
	const deductQuery = `
		UPDATE honey.user_honey SET honey = honey - $1 WHERE user_id = $2
	`
	const insertQuery = `
		INSERT INTO honey.active_challenges (creator_user_id, amount)
		VALUES ($1, $2)
		RETURNING id
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.ActiveChallenge{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance int64
	if err := tx.QueryRow(ctx, lockHoneyQuery, creatorUserID).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ActiveChallenge{}, domain.ErrUserHoneyNotFound
		}
		return domain.ActiveChallenge{}, fmt.Errorf("select honey: %w", err)
	}

	if balance < amount {
		return domain.ActiveChallenge{}, domain.ErrNotEnoughHoney
	}

	if _, err := tx.Exec(ctx, deductQuery, amount, creatorUserID); err != nil {
		return domain.ActiveChallenge{}, fmt.Errorf("deduct honey: %w", err)
	}

	var id int
	if err := tx.QueryRow(ctx, insertQuery, creatorUserID, amount).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case uniqueViolationCode:
				return domain.ActiveChallenge{}, domain.ErrChallengeAlreadyExists
			case checkViolationCode:
				return domain.ActiveChallenge{}, domain.ErrWrongAmount
			}
		}
		return domain.ActiveChallenge{}, fmt.Errorf("insert active challenge: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ActiveChallenge{}, fmt.Errorf("commit transaction: %w", err)
	}

	return domain.NewActiveChallenge(id, creatorUserID, "", amount), nil
}
