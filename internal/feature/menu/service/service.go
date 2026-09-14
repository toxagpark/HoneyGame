package menu_service

import (
	"context"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
	fight_service "github.com/toxagpark/HoneyGame/internal/feature/fight/service"
)

// BetOption — вариант ставки в меню: процент от баланса и округлённая сумма.
type BetOption struct {
	Percent int
	Amount  int64
}

// betPercents — проценты баланса, из которых игрок выбирает ставку.
var betPercents = []int{10, 15, 20, 25}

// minBetBalance — минимальный баланс, чтобы самый маленький вариант (10%)
// давал хотя бы 1 мёд ставки.
const minBetBalance = 5

type Service struct {
	users usersService
	fight *fight_service.Service
}

type usersService interface {
	GetUser(
		ctx context.Context,
		tgChatID int64,
	) (domain.User, error)
}

func NewService(
	users usersService,
	fight *fight_service.Service,
) *Service {
	return &Service{
		users: users,
		fight: fight,
	}
}

// Menu возвращает профиль игрока для главного экрана меню.
func (s *Service) Profile(
	ctx context.Context,
	tgChatID int64,
) (domain.User, error) {
	return s.users.GetUser(ctx, tgChatID)
}

// BetOptions строит варианты ставок от текущего баланса: 10/15/20/25%,
// округлённые вниз до целых. Пустой слайс — баланс слишком мал для боя.
// Одинаковые суммы после округления схлопываются: остаётся меньший процент.
func (s *Service) BetOptions(
	ctx context.Context,
	honey int64,
) []BetOption {
	if honey < minBetBalance {
		return nil
	}

	options := make([]BetOption, 0, len(betPercents))
	for _, pct := range betPercents {
		amount := honey * int64(pct) / 100
		if amount < 1 {
			continue
		}
		// Ставка в мёде уже встречалась (округление съело разницу) — пропускаем.
		if len(options) > 0 && options[len(options)-1].Amount == amount {
			continue
		}
		options = append(options, BetOption{Percent: pct, Amount: amount})
	}

	return options
}

// CreateChallengeFromMenu создаёт вызов на pct% от текущего баланса игрока.
func (s *Service) CreateChallengeFromMenu(
	ctx context.Context,
	tgChatID int64,
	pct int,
) (domain.ActiveChallenge, error) {
	user, err := s.users.GetUser(ctx, tgChatID)
	if err != nil {
		return domain.ActiveChallenge{}, err
	}

	amount := user.Honey * int64(pct) / 100
	if amount < 1 {
		return domain.ActiveChallenge{}, domain.ErrNotEnoughHoney
	}

	return s.fight.CreateChallenge(ctx, tgChatID, amount)
}

// Fights возвращает открытые вызовы других игроков.
func (s *Service) Fights(
	ctx context.Context,
	tgChatID int64,
) ([]domain.ActiveChallenge, error) {
	return s.fight.GetChallenges(ctx, tgChatID)
}

// MyChallenge возвращает активный вызов игрока; ErrChallengeNotFound, если его нет.
func (s *Service) MyChallenge(
	ctx context.Context,
	tgChatID int64,
) (domain.ActiveChallenge, error) {
	return s.fight.MyChallenge(ctx, tgChatID)
}
