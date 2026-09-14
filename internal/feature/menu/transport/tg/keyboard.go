package menu_tg_transport

import (
	"errors"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func backRow() [][]telego.InlineKeyboardButton {
	return [][]telego.InlineKeyboardButton{
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData(callbackMenuBack)),
	}
}

func tryStartText(err error) string {
	switch {
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrUserHoneyNotFound):
		return "Ты ещё не зарегистрирован 🐻 Напиши /start"
	default:
		return "Ошибка( Попробуйте прописать /start"
	}
}

func createChallengeErrorText(err error) string {
	switch {
	case errors.Is(err, domain.ErrChallengeAlreadyExists):
		return "У тебя уже есть активный вызов! Отзови его или дождись боя 🐻"
	case errors.Is(err, domain.ErrNotEnoughHoney):
		return "Не хватает мёда на такую ставку 🍯🐻"
	default:
		return tryStartText(err)
	}
}
