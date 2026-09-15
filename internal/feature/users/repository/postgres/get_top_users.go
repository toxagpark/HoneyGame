package users_pg_repo

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// GetTopUsers возвращает игроков с наибольшим балансом мёда, максимум limit.
func (r *Repository) GetTopUsers(
	ctx context.Context,
	limit int,
) ([]domain.User, error) {
	const query = `
		SELECT u.id, u.tg_user_id, u.user_name, h.honey
		FROM honey.users u
		JOIN honey.user_honey h ON h.user_id = u.id
		ORDER BY h.honey DESC, u.id
		LIMIT $1
	`

	rows, err := r.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var id int
		var tgUserID int64
		var userName string
		var honey int64
		if err := rows.Scan(&id, &tgUserID, &userName, &honey); err != nil {
			return nil, err
		}

		users = append(
			users,
			domain.NewUser(id, tgUserID, userName, honey),
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
