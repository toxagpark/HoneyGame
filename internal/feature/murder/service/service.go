package murder_service

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type Service struct {
	repo  repository
	users usersService
}

type repository interface {
	Rob(
		ctx context.Context,
		userID int,
		hives int,
		amount int64,
	) (domain.MurderResult, error)
}

// usersService — узкий интерфейс к фиче users: перевод tg_user_id во внутренний ID.
type usersService interface {
	GetUser(
		ctx context.Context,
		tgUserID int64,
	) (domain.User, error)
}

func NewService(
	r repository,
	users usersService,
) *Service {
	return &Service{
		repo:  r,
		users: users,
	}
}
