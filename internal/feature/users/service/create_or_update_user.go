package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (s *Service) CreateOrUpdateUser(
	ctx context.Context,
	newUser domain.User,
) (domain.User, error) {
	user, err := s.CreateUser(ctx, newUser)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		return domain.User{}, fmt.Errorf("create or update user: %w", err)
	}

	user, err = s.UpdateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
