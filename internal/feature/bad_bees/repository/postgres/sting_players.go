package bad_bees_pg_repo

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// StingPlayers отнимает у каждого медведя его штраф (параллельные массивы:
// i-му игроку — i-й штраф, уже посчитанный от его баланса) и возвращает
// имена с новыми балансами.
func (r *Repository) StingPlayers(
	ctx context.Context,
	userIDs []int,
	stings []int64,
) ([]domain.User, error) {
	const query = `
		WITH stung AS (
			UPDATE honey.user_honey h
			SET honey = GREATEST(h.honey - s.sting, 0)
			FROM (
				SELECT id, sting
				FROM unnest($1::int[], $2::bigint[]) AS t(id, sting)
			) s
			WHERE h.user_id = s.id
			RETURNING h.user_id, h.honey
		)
		SELECT st.user_id, u.user_name, st.honey
		FROM stung st
		JOIN honey.users u ON u.id = st.user_id
		ORDER BY st.honey DESC
	`

	rows, err := r.Pool.Query(ctx, query, userIDs, stings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var id int
		var userName string
		var honey int64
		if err := rows.Scan(&id, &userName, &honey); err != nil {
			return nil, err
		}

		users = append(users, domain.User{
			ID:       id,
			UserName: userName,
			Honey:    honey,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
