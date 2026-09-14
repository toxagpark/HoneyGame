package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (s *Service) CreateUser(
	ctx context.Context,
	newUser domain.User,
) (domain.User, error) {
	if newUser.UserName == "" {
		newUser.UserName = "Неизвестный медведь"
	}

	user, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return domain.User{}, err
		}
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
