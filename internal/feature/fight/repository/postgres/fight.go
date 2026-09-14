package fight_pg_repo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type FightResult struct {
	WinnerUserID int
	LoserUserID  int
	Amount       int64
}

// Fight проводит бой. Ставка создателя уже в банке (списана при создании вызова),
// поэтому здесь списываем только у принимающего и начисляем победителю x2.
func (r *Repository) Fight(
	ctx context.Context,
	creatorUserID int,
	acceptorUserID int,
	amount int64,
) (FightResult, error) {
	const challengeQuery = `
		DELETE FROM honey.active_challenges
		WHERE creator_user_id = $1 AND amount = $2
		RETURNING id
	`
	const honeyQuery = `
		SELECT honey FROM honey.user_honey WHERE user_id = $1 FOR UPDATE
	`
	const acceptorUpdateQuery = `
		UPDATE honey.user_honey SET honey = honey - $1 WHERE user_id = $2
	`
	const winnerUpdateQuery = `
		UPDATE honey.user_honey SET honey = honey + $1 WHERE user_id = $2
	`
	const logQuery = `
		INSERT INTO honey.challenges (amount, winner_user_id, loser_user_id, completed_at)
		VALUES ($1, $2, $3, NOW())
	`

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return FightResult{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Блокируем строку вызова: только один принимающий сможет его забрать.
	var challengeID int
	if err := tx.QueryRow(ctx, challengeQuery, creatorUserID, amount).Scan(&challengeID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FightResult{}, domain.ErrChallengeNotFound
		}
		return FightResult{}, fmt.Errorf("delete active challenge: %w", err)
	}

	// FOR UPDATE в одном порядке (по возрастанию user_id) защищает от дедлока,
	// когда одни и те же медведи дерутся друг с другом встречными вызовами.
	firstID, secondID := creatorUserID, acceptorUserID
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	var firstHoney, secondHoney int64
	if err := tx.QueryRow(ctx, honeyQuery, firstID).Scan(&firstHoney); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FightResult{}, domain.ErrUserHoneyNotFound
		}
		return FightResult{}, fmt.Errorf("select first honey: %w", err)
	}
	if err := tx.QueryRow(ctx, honeyQuery, secondID).Scan(&secondHoney); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return FightResult{}, domain.ErrUserHoneyNotFound
		}
		return FightResult{}, fmt.Errorf("select second honey: %w", err)
	}

	// first/second — порядок блокировки, а не creator/acceptor.
	// Баланс создателя уже в банке, проверяем только принимающего.
	acceptorBalance := firstHoney
	if acceptorUserID == secondID {
		acceptorBalance = secondHoney
	}
	if acceptorBalance < amount {
		return FightResult{}, domain.ErrNotEnoughHoney
	}

	// 50/50
	creatorWins := rand.Intn(2) == 0

	var winnerID, loserID int
	if creatorWins {
		winnerID, loserID = creatorUserID, acceptorUserID
	} else {
		winnerID, loserID = acceptorUserID, creatorUserID
	}

	if _, err := tx.Exec(ctx, acceptorUpdateQuery, amount, acceptorUserID); err != nil {
		return FightResult{}, fmt.Errorf("update acceptor honey: %w", err)
	}
	if _, err := tx.Exec(ctx, winnerUpdateQuery, amount*2, winnerID); err != nil {
		return FightResult{}, fmt.Errorf("update winner honey: %w", err)
	}

	if _, err := tx.Exec(ctx, logQuery, amount, winnerID, loserID); err != nil {
		return FightResult{}, fmt.Errorf("insert challenge log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return FightResult{}, fmt.Errorf("commit transaction: %w", err)
	}

	// Amount — чистый выигрыш/проигрыш: ставка, а не x2 из банка.
	return FightResult{
		WinnerUserID: winnerID,
		LoserUserID:  loserID,
		Amount:       amount,
	}, nil
}
