package fight_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// MyChallenge возвращает активный вызов игрока; ErrChallengeNotFound, если его нет.
func (s *Service) MyChallenge(
	ctx context.Context,
	tgChatID int64,
) (domain.ActiveChallenge, error) {
	user, err := s.users.GetUser(ctx, tgChatID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to get user: %w", err)
	}

	challenge, err := s.repo.GetChallengeByCreator(ctx, user.ID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			return domain.ActiveChallenge{}, err
		}
		return domain.ActiveChallenge{}, fmt.Errorf("failed to get own challenge: %w", err)
	}

	return challenge, nil
}
