package bad_bees_pg_repo

import (
	"context"
	"time"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// GetIdlePlayers возвращает медведей без боевой активности за период since:
// ни боёв в журнале challenges, ни свежего вызова в active_challenges —
// непринятый вызов тоже активность, трусы не приняли, а храбрец не АФК.
// Пустой баланс не отнимешь — таких пчёлы игнорируют.
func (r *Repository) GetIdlePlayers(
	ctx context.Context,
	since time.Time,
) ([]domain.User, error) {
	const query = `
		SELECT u.id, u.tg_user_id, u.user_name, h.honey
		FROM honey.users u
		JOIN honey.user_honey h ON h.user_id = u.id
		WHERE h.honey > 0
		  AND NOT EXISTS (
			SELECT 1
			FROM honey.challenges c
			WHERE (c.winner_user_id = u.id OR c.loser_user_id = u.id)
			  AND c.completed_at >= $1
		  )
		  AND NOT EXISTS (
			SELECT 1
			FROM honey.active_challenges a
			WHERE a.creator_user_id = u.id
			  AND a.created_at >= $1
		  )
	`

	rows, err := r.Pool.Query(ctx, query, since)
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
