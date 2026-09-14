package fight_pg_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// DeleteChallenge удаляет вызов и возвращает создателю его ставку одной транзакцией.
func (r *Repository) DeleteChallenge(
	ctx context.Context,
	challengeID int,
	creatorUserID int,
) (domain.ActiveChallenge, error) {
	const deleteQuery = `
		DELETE FROM honey.active_challenges
		WHERE id = $1 AND creator_user_id = $2
		RETURNING amount
	`
	const refundQuery = `
		UPDATE honey.user_honey SET honey = honey + $1 WHERE user_id = $2
		RETURNING honey
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.ActiveChallenge{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var amount int64
	if err := tx.QueryRow(ctx, deleteQuery, challengeID, creatorUserID).Scan(&amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ActiveChallenge{}, domain.ErrChallengeNotFound
		}
		return domain.ActiveChallenge{}, fmt.Errorf("delete active challenge: %w", err)
	}

	var refunded int64
	if err := tx.QueryRow(ctx, refundQuery, amount, creatorUserID).Scan(&refunded); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ActiveChallenge{}, domain.ErrUserHoneyNotFound
		}
		return domain.ActiveChallenge{}, fmt.Errorf("refund honey: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ActiveChallenge{}, fmt.Errorf("commit transaction: %w", err)
	}

	return domain.NewActiveChallenge(challengeID, creatorUserID, "", amount), nil
}
