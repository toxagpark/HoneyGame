package murder_pg_repo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// robTimeout ограничивает транзакцию: если не уложились — откат без списаний.
const robTimeout = 5 * time.Second

// Rob проводит ограбление улья одной транзакцией: блокирует баланс, проверяет
// ставку и разыгрывает улей с мёдом. Проигрыш = −ставка,
// победа = +ставка × количество ульев.
func (r *Repository) Rob(
	ctx context.Context,
	userID int,
	hives int,
	amount int64,
) (domain.MurderResult, error) {
	ctx, cancel := context.WithTimeout(ctx, robTimeout)
	defer cancel()

	const honeyQuery = `
		SELECT honey FROM honey.user_honey WHERE user_id = $1 FOR UPDATE
	`
	const updateQuery = `
		UPDATE honey.user_honey SET honey = honey + $1 WHERE user_id = $2 RETURNING honey
	`
	const logQuery = `
		INSERT INTO honey.murders (user_id, hives, amount, honey)
		VALUES ($1, $2, $3, $4)
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.MurderResult{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var balance int64
	if err := tx.QueryRow(ctx, honeyQuery, userID).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MurderResult{}, domain.ErrUserHoneyNotFound
		}
		return domain.MurderResult{}, fmt.Errorf("select honey: %w", err)
	}

	if balance < amount {
		return domain.MurderResult{}, domain.ErrNotEnoughHoney
	}

	// Мёд лежит в одном случайном улье, лапа тычет в свой случайный:
	// победа — только когда они совпали.
	pawHive := rand.Intn(hives) + 1
	honeyHive := rand.Intn(hives) + 1

	// Honey — чистый итог ограбления для журнала и чата.
	honey := -amount
	if pawHive == honeyHive {
		honey = amount * int64(hives)
	}

	var newHoney int64
	if err := tx.QueryRow(ctx, updateQuery, honey, userID).Scan(&newHoney); err != nil {
		return domain.MurderResult{}, fmt.Errorf("update honey: %w", err)
	}

	if _, err := tx.Exec(ctx, logQuery, userID, hives, amount, honey); err != nil {
		return domain.MurderResult{}, fmt.Errorf("insert murder log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.MurderResult{}, fmt.Errorf("commit transaction: %w", err)
	}

	return domain.NewMurderResult(pawHive, honeyHive, honey, newHoney), nil
}
