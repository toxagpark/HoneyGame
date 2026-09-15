package users_service

import (
	"context"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (s *Service) GetUser(
	ctx context.Context,
	tgUserID int64,
) (domain.User, error) {
	user, err := s.repo.GetUser(
		ctx,
		tgUserID,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
