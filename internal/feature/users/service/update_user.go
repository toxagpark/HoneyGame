package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (s *Service) UpdateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if user.UserName == "" {
		user.UserName = "Неизвестный медведь"
	}

	updatedUser, err := s.repo.UpdateUser(
		ctx,
		user,
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, err
		}
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}
