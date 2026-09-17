package murder_service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

const (
	minHives = 3
	maxHives = 10
	// minStake — улей дешевле одного мёда не грабят.
	minStake int64 = 1
)

// RobResult — исход ограбления с точки зрения чата: имя игрока из БД,
// куда ткнула лапа, где был мёд, чистый выигрыш (отрицательный — проигрыш)
// и новый баланс.
type RobResult struct {
	UserName  string
	PawHive   int
	HoneyHive int
	Honey     int64
	NewHoney  int64
}

// HivesBounds — допустимые границы количества ульев: нужны транспорту для текста правил.
func HivesBounds() (int, int) {
	return minHives, maxHives
}

// Rob проводит ограбление улья (по tg_user_id): игрок ставит amount мёда,
// мишка тычет лапу в один из hives ульев, мёд лежит в одном из них.
// Угадал — ставка × количество ульев, не угадал — ставка остаётся пчёлам.
func (s *Service) Rob(
	ctx context.Context,
	tgUserID int64,
	hives int,
	amount int64,
) (RobResult, error) {
	if hives < minHives || hives > maxHives {
		return RobResult{}, domain.ErrWrongHives
	}
	if amount < minStake || amount > math.MaxInt64/int64(hives) {
		// Верхняя граница — чтобы ставка × ульи не переполнилась.
		return RobResult{}, domain.ErrWrongAmount
	}

	user, err := s.users.GetUser(ctx, tgUserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return RobResult{}, err
		}
		return RobResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	result, err := s.repo.Rob(ctx, user.ID, hives, amount)
	if err != nil {
		if errors.Is(err, domain.ErrNotEnoughHoney) ||
			errors.Is(err, domain.ErrUserHoneyNotFound) {
			return RobResult{}, err
		}
		return RobResult{}, fmt.Errorf("failed to rob: %w", err)
	}

	return RobResult{
		UserName:  user.UserName,
		PawHive:   result.PawHive,
		HoneyHive: result.HoneyHive,
		Honey:     result.Honey,
		NewHoney:  result.NewHoney,
	}, nil
}
