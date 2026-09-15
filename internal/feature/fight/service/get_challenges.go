package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// GetChallenges возвращает открытые вызовы других игроков (свои исключены).
func (s *Service) GetChallenges(
	ctx context.Context,
	tgUserID int64,
) ([]domain.ActiveChallenge, error) {
	user, err := s.users.GetUser(ctx, tgUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	challenges, err := s.repo.GetChallenges(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenges: %w", err)
	}

	return challenges, nil
}
