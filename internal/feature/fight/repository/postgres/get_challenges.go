package fight_pg_repo

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// GetChallenges возвращает открытые вызовы других игроков, кроме excludeUserID.
func (r *Repository) GetChallenges(
	ctx context.Context,
	excludeUserID int,
) ([]domain.ActiveChallenge, error) {
	const query = `
		SELECT c.id, c.creator_user_id, u.user_name, c.amount
		FROM honey.active_challenges c
		JOIN honey.users u ON u.id = c.creator_user_id
		WHERE c.creator_user_id != $1
		ORDER BY c.created_at
	`

	rows, err := r.Pool.Query(ctx, query, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := []domain.ActiveChallenge{}
	for rows.Next() {
		var id, creatorUserID int
		var creatorName string
		var amount int64
		if err := rows.Scan(&id, &creatorUserID, &creatorName, &amount); err != nil {
			return nil, err
		}

		challenges = append(
			challenges,
			domain.NewActiveChallenge(id, creatorUserID, creatorName, amount),
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}
