package bad_bees_service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
	"github.com/toxagpark/HoneyGame/internal/feature/bad_bees"
)

type Service struct {
	repo          repository
	rand          *rand.Rand
	recentPeriod  time.Duration
	maxStingRatio float64
}

// Sting — итог одного укуса: сколько мёда отнято и какой баланс остался.
type Sting struct {
	UserName string
	Honey    int64
	NewHoney int64
}

type repository interface {
	GetIdlePlayers(
		ctx context.Context,
		since time.Time,
	) ([]domain.User, error)

	StingPlayers(
		ctx context.Context,
		userIDs []int,
		stings []int64,
	) ([]domain.User, error)
}

func NewService(
	repo repository,
	cfg *bad_bees.Config,
) *Service {
	return &Service{
		repo:          repo,
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		recentPeriod:  cfg.RECENT_PERIOD,
		maxStingRatio: float64(cfg.MAX_STING_PERCENT) / 100,
	}
}

// StingIdle находит медведей без боевой активности за последнее время
// (без боёв и свежих вызовов) и отнимает у каждого случайный процент мёда
// (от 1% до максимума из конфига). Возвращает итоги по всем укушенным.
func (s *Service) StingIdle(
	ctx context.Context,
) ([]Sting, error) {
	since := time.Now().Add(-s.recentPeriod)

	idle, err := s.repo.GetIdlePlayers(ctx, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get idle players: %w", err)
	}
	if len(idle) == 0 {
		return []Sting{}, nil
	}

	// i-му игроку — i-й укус: случайный процент от его текущего баланса.
	userIDs := make([]int, len(idle))
	stings := make([]int64, len(idle))
	for i, u := range idle {
		userIDs[i] = u.ID
		percent := 1 + s.rand.Intn(int(s.maxStingRatio*100))
		stings[i] = max(int64(float64(u.Honey)*float64(percent)/100), 1)
	}

	stung, err := s.repo.StingPlayers(ctx, userIDs, stings)
	if err != nil {
		return nil, fmt.Errorf("failed to sting players: %w", err)
	}

	// Репозиторий возвращает игроков в своём порядке, склеиваем по ID.
	stingByID := make(map[int]int64, len(userIDs))
	for i, id := range userIDs {
		stingByID[id] = stings[i]
	}

	results := make([]Sting, len(stung))
	for i, u := range stung {
		results[i] = Sting{
			UserName: u.UserName,
			Honey:    stingByID[u.ID],
			NewHoney: u.Honey,
		}
	}

	return results, nil
}
