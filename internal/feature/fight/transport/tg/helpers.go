package fight_tg_transport

import (
	"errors"
	"strconv"
	"strings"

	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func parseAmount(text string) (int64, error) {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) < 2 {
		return 0, errors.New("amount is required")
	}
	return strconv.ParseInt(parts[1], 10, 64)
}

func tryStartText(err error) string {
	switch {
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrUserHoneyNotFound):
		return "Ты ещё не зарегистрирован 🐻 Напиши /start"
	default:
		return "Ошибка( Попробуйте прописать /start"
	}
}
