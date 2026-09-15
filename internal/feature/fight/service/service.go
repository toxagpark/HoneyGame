package fight_service

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type Service struct {
	repo  repository
	users usersService
}

type repository interface {
	CreateChallenge(
		ctx context.Context,
		creatorUserID int,
		amount int64,
	) (domain.ActiveChallenge, error)

	GetChallenge(
		ctx context.Context,
		challengeID int,
	) (domain.ActiveChallenge, error)

	GetChallengeByCreator(
		ctx context.Context,
		creatorUserID int,
	) (domain.ActiveChallenge, error)

	GetChallenges(
		ctx context.Context,
		excludeUserID int,
	) ([]domain.ActiveChallenge, error)

	DeleteChallenge(
		ctx context.Context,
		challengeID int,
		creatorUserID int,
	) (domain.ActiveChallenge, error)

	Fight(
		ctx context.Context,
		creatorUserID int,
		acceptorUserID int,
		amount int64,
	) (domain.FightResult, error)
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
