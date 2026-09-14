package users_pg_repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (r *Repository) UpdateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	const query = `
		UPDATE honey.users u
		SET user_name = $2
		WHERE u.tg_chat_id = $1
		RETURNING u.id, u.tg_chat_id, u.user_name,
			(SELECT h.honey FROM honey.user_honey h WHERE h.user_id = u.id)
	`

	row := r.Pool.QueryRow(
		ctx,
		query,
		user.TgChatID,
		user.UserName,
	)

	var id int
	var tgChatId int64
	var userName string
	var honey *int64
	if err := row.Scan(&id, &tgChatId, &userName, &honey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}

	if honey == nil {
		return domain.User{}, domain.ErrUserHoneyNotFound
	}

	updatedUser := domain.NewUser(
		id,
		tgChatId,
		userName,
		*honey,
	)

	return updatedUser, nil
}
