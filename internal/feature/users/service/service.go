package users_service

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type Service struct {
	repo repository
}

type repository interface {
	CreateUser(
		ctx context.Context,
		newUser domain.User,
	) (domain.User, error)

	GetUser(
		ctx context.Context,
		tgUserID int64,
	) (domain.User, error)

	GetUserByID(
		ctx context.Context,
		userID int,
	) (domain.User, error)

	UpdateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetTopUsers(
		ctx context.Context,
		limit int,
	) ([]domain.User, error)
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}
