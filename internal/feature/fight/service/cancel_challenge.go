package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// CancelChallenge отзывает свой вызов и возвращает ставку.
func (s *Service) CancelChallenge(
	ctx context.Context,
	tgUserID int64,
	challengeID int,
) (domain.ActiveChallenge, error) {
	user, err := s.users.GetUser(ctx, tgUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to get user: %w", err)
	}

	challenge, err := s.repo.GetChallenge(ctx, challengeID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to get challenge: %w", err)
	}

	if challenge.CreatorUserID != user.ID {
		return domain.ActiveChallenge{}, domain.ErrChallengeNotOwner
	}

	deleted, err := s.repo.DeleteChallenge(ctx, challengeID, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to cancel challenge: %w", err)
	}

	return deleted, nil
}
