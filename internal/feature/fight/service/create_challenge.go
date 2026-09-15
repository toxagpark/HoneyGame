package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

const minChallengeAmount = 1

// CreateChallenge создаёт вызов от игрока (по tg_user_id) со списанием ставки.
func (s *Service) CreateChallenge(
	ctx context.Context,
	tgUserID int64,
	amount int64,
) (domain.ActiveChallenge, error) {
	if amount < minChallengeAmount {
		return domain.ActiveChallenge{}, domain.ErrWrongAmount
	}

	user, err := s.users.GetUser(ctx, tgUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to get user: %w", err)
	}

	challenge, err := s.repo.CreateChallenge(ctx, user.ID, amount)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeAlreadyExists) ||
			errors.Is(err, domain.ErrWrongAmount) ||
			errors.Is(err, domain.ErrNotEnoughHoney) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to create challenge: %w", err)
	}

	challenge.CreatorName = user.UserName

	return challenge, nil
}
