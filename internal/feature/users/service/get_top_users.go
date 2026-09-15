package users_service

import (
	"context"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// topUsersLimit — сколько игроков показывать в /top.
const topUsersLimit = 10

// GetTopUsers возвращает топ игроков по балансу мёда (по убыванию).
func (s *Service) GetTopUsers(
	ctx context.Context,
) ([]domain.User, error) {
	users, err := s.repo.GetTopUsers(ctx, topUsersLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top users: %w", err)
	}

	return users, nil
}
