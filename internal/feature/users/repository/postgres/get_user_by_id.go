package users_pg_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (r *Repository) GetUserByID(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	const query = `
		SELECT u.id, u.tg_user_id, u.user_name, h.honey
		FROM honey.users u
		JOIN honey.user_honey h ON h.user_id = u.id
		WHERE u.id = $1
	`

	row := r.Pool.QueryRow(
		ctx,
		query,
		userID,
	)

	var id int
	var dbTgUserID int64
	var userName string
	var honey int64
	if err := row.Scan(&id, &dbTgUserID, &userName, &honey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	user := domain.NewUser(
		id,
		dbTgUserID,
		userName,
		honey,
	)

	return user, nil
}
