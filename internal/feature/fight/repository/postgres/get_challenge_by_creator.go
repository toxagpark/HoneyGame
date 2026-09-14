package fight_pg_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// GetChallengeByCreator возвращает активный вызов игрока, если он есть.
func (r *Repository) GetChallengeByCreator(
	ctx context.Context,
	creatorUserID int,
) (domain.ActiveChallenge, error) {
	const query = `
		SELECT c.id, c.creator_user_id, u.user_name, c.amount
		FROM honey.active_challenges c
		JOIN honey.users u ON u.id = c.creator_user_id
		WHERE c.creator_user_id = $1
	`

	row := r.Pool.QueryRow(ctx, query, creatorUserID)

	var id, creatorID int
	var creatorName string
	var amount int64
	if err := row.Scan(&id, &creatorID, &creatorName, &amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ActiveChallenge{}, domain.ErrChallengeNotFound
		}
		return domain.ActiveChallenge{}, err
	}

	return domain.NewActiveChallenge(id, creatorID, creatorName, amount), nil
}
