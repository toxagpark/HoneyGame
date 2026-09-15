package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// AcceptChallengeResult — исход боя с точки зрения чата: имена победителя
// и проигравшего из БД и ставка.
type AcceptChallengeResult struct {
	WinnerName string
	LoserName  string
	Amount     int64
}

// ValidateAccept проверяет, можно ли принять вызов. Нужен транспорту,
// чтобы отсечь ошибки (в т.ч. нехватку мёда у принимающего) до анимации боя.
func (s *Service) ValidateAccept(
	ctx context.Context,
	tgUserID int64,
	challengeID int,
) error {
	user, err := s.users.GetUser(ctx, tgUserID)
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

// AcceptChallenge принимает вызов (по tg_user_id), проводит бой и возвращает исход.
func (s *Service) AcceptChallenge(
	ctx context.Context,
	tgUserID int64,
	challengeID int,
) (AcceptChallengeResult, error) {
	user, err := s.users.GetUser(ctx, tgUserID)
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

	// Имена обоих участников — из БД: у создателя имя пришло с вызовом,
	// у принявшего — из профиля (для безюзернеймных оно сгенерировано при /start).
	winnerName := user.UserName
	loserName := challenge.CreatorName
	if result.WinnerUserID == challenge.CreatorUserID {
		winnerName = challenge.CreatorName
		loserName = user.UserName
	}

	return AcceptChallengeResult{
		WinnerName: winnerName,
		LoserName:  loserName,
		Amount:     result.Amount,
	}, nil
}
