package honey_day_service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

// minHoneyDayGift и maxHoneyDayGift — границы личного подарка медового дня.
const (
	minHoneyDayGift = 5
	maxHoneyDayGift = 25
)

type Service struct {
	repo repository
	rand *rand.Rand
}

// HoneyDayGift — персональный подарок одного медового дня.
type HoneyDayGift struct {
	UserName string
	Honey    int64
	NewHoney int64
}

type repository interface {
	GetPlayerIDs(
		ctx context.Context,
	) ([]int, error)

	AddHoneyToPlayers(
		ctx context.Context,
		userIDs []int,
		gifts []int64,
	) ([]domain.User, error)
}

func NewService(
	repo repository,
) *Service {
	return &Service{
		repo: repo,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GiveHoneyDay раздаёт каждому игроку свой случайный подарок [5, 25]
// и возвращает итоги по всем участникам.
func (s *Service) GiveHoneyDay(
	ctx context.Context,
) ([]HoneyDayGift, error) {
	userIDs, err := s.repo.GetPlayerIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get player ids: %w", err)
	}
	if len(userIDs) == 0 {
		return []HoneyDayGift{}, nil
	}

	gifts := make([]int64, len(userIDs))
	for i := range gifts {
		gifts[i] = minHoneyDayGift + s.rand.Int63n(maxHoneyDayGift-minHoneyDayGift+1)
	}

	users, err := s.repo.AddHoneyToPlayers(ctx, userIDs, gifts)
	if err != nil {
		return nil, fmt.Errorf("failed to add honey to players: %w", err)
	}

	// Репозиторий возвращает игроков в своём порядке, склеиваем подарки по ID.
	giftByID := make(map[int]int64, len(userIDs))
	for i, id := range userIDs {
		giftByID[id] = gifts[i]
	}

	giftResults := make([]HoneyDayGift, len(users))
	for i, u := range users {
		giftResults[i] = HoneyDayGift{
			UserName: u.UserName,
			Honey:    giftByID[u.ID],
			NewHoney: u.Honey,
		}
	}

	return giftResults, nil
}
