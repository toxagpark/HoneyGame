package users_pg_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

const uniqueViolationCode = "23505"

// registrationHoney — стартовый бонус мёда за регистрацию.
const registrationHoney = 15

func (r *Repository) CreateUser(
	ctx context.Context,
	newUser domain.User,
) (domain.User, error) {
	const query = `
		INSERT INTO honey.users (tg_chat_id, user_name)
		VALUES ($1, $2)
		RETURNING id, tg_chat_id, user_name
	`
	const honeyQuery = `
		INSERT INTO honey.user_honey (user_id, honey)
		VALUES ($1, $2)
		RETURNING honey
	`
	const ensureHoneyQuery = `
		INSERT INTO honey.user_honey (user_id)
		SELECT id FROM honey.users WHERE tg_chat_id = $1
		ON CONFLICT (user_id) DO NOTHING
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(
		ctx,
		query,
		newUser.TgChatID,
		newUser.UserName,
	)

	var id int
	var tgChatId int64
	var userName string
	if err := row.Scan(&id, &tgChatId, &userName); err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != uniqueViolationCode {
			return domain.User{}, fmt.Errorf("insert user: %w", err)
		}

		// После ошибки транзакция прервана (25P02), работать в ней дальше нельзя:
		// откатываем и чиним user_honey отдельным автокоммитным запросом.
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return domain.User{}, fmt.Errorf("rollback transaction: %w", rbErr)
		}
		if _, err := r.Pool.Exec(ctx, ensureHoneyQuery, newUser.TgChatID); err != nil {
			return domain.User{}, fmt.Errorf("ensure user honey: %w", err)
		}

		return domain.User{}, domain.ErrUserAlreadyExists
	}

	honeyRow := tx.QueryRow(
		ctx,
		honeyQuery,
		id,
		registrationHoney,
	)

	var honey int64
	if err := honeyRow.Scan(&honey); err != nil {
		return domain.User{}, fmt.Errorf("insert user honey: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit transaction: %w", err)
	}

	user := domain.NewUser(
		id,
		tgChatId,
		userName,
		honey,
	)

	return user, nil
}
