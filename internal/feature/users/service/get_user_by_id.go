package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (s *Service) GetUserByID(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	user, err := s.repo.GetUserByID(
		ctx,
		userID,
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, err
		}
		return domain.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
