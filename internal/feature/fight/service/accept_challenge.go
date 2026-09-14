package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// AcceptChallengeResult — исход боя с точки зрения принявшего вызов.
type AcceptChallengeResult struct {
	AcceptorWon      bool
	Amount           int64
	OpponentTgChatID int64
}

// ValidateAccept проверяет, можно ли принять вызов. Нужен транспорту,
// чтобы отсечь ошибки (в т.ч. нехватку мёда у принимающего) до анимации боя.
func (s *Service) ValidateAccept(
	ctx context.Context,
	tgChatID int64,
	challengeID int,
) error {
	user, err := s.users.GetUser(ctx, tgChatID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return err
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	challenge, err := s.repo.GetChallenge(ctx, challengeID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			return err
		}
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge.CreatorUserID == user.ID {
		return domain.ErrSelfChallenge
	}

	// Ставка принимающего списывается в момент боя: проверяем её заранее,
	// чтобы не «играть» бой, который заведомо не состоится.
	if user.Honey < challenge.Amount {
		return domain.ErrNotEnoughHoney
	}

	return nil
}

// AcceptChallenge принимает вызов (по tg_chat_id), проводит бой и возвращает исход.
func (s *Service) AcceptChallenge(
	ctx context.Context,
	tgChatID int64,
	challengeID int,
) (AcceptChallengeResult, error) {
	user, err := s.users.GetUser(ctx, tgChatID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return AcceptChallengeResult{}, err
		}
		return AcceptChallengeResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	challenge, err := s.repo.GetChallenge(ctx, challengeID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			return AcceptChallengeResult{}, err
		}
		return AcceptChallengeResult{}, fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge.CreatorUserID == user.ID {
		return AcceptChallengeResult{}, domain.ErrSelfChallenge
	}

	result, err := s.repo.Fight(ctx, challenge.CreatorUserID, user.ID, challenge.Amount)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) ||
			errors.Is(err, domain.ErrNotEnoughHoney) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return AcceptChallengeResult{}, err
		}
		return AcceptChallengeResult{}, fmt.Errorf("failed to fight: %w", err)
	}

	opponent, err := s.users.GetUserByID(ctx, challenge.CreatorUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return AcceptChallengeResult{}, err
		}
		return AcceptChallengeResult{}, fmt.Errorf("failed to get opponent: %w", err)
	}

	return AcceptChallengeResult{
		AcceptorWon:      result.WinnerUserID == user.ID,
		Amount:           result.Amount,
		OpponentTgChatID: opponent.TgChatID,
	}, nil
}
