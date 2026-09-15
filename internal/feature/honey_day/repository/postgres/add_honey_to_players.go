package honey_day_pg_repo

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// AddHoneyToPlayers начисляет каждому игроку его сумму из gifts (параллельные
// массивы: i-му игроку — i-й подарок) и возвращает имена с новыми балансами.
func (r *Repository) AddHoneyToPlayers(
	ctx context.Context,
	userIDs []int,
	gifts []int64,
) ([]domain.User, error) {
	const query = `
		WITH paid AS (
			UPDATE honey.user_honey h
			SET honey = honey + g.gift
			FROM (
				SELECT id, gift
				FROM unnest($1::int[], $2::bigint[]) AS t(id, gift)
			) g
			WHERE h.user_id = g.id
			RETURNING h.user_id, h.honey
		)
		SELECT p.user_id, u.user_name, p.honey
		FROM paid p
		JOIN honey.users u ON u.id = p.user_id
		ORDER BY p.honey DESC
	`

	rows, err := r.Pool.Query(ctx, query, userIDs, gifts)
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

// GetPlayerIDs возвращает ID всех игроков — кому воркер раздаёт мёд.
func (r *Repository) GetPlayerIDs(
	ctx context.Context,
) ([]int, error) {
	const query = `SELECT id FROM honey.users`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
