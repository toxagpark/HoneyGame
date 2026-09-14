package users_tg_transport

import (
	"context"
	"errors"
	"fmt"
	"log"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (h *Handler) HandleGetUser(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	user, err := h.service.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		var text string
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			text = "Ты ещё не зарегистрирован 🐻 Напиши /start"
		case errors.Is(err, domain.ErrUserHoneyNotFound):
			text = "Ой, твой мёд потерялся 🐻 Напиши /start"
		default:
			text = "Ошибка получения пользователя("
			log.Println(err)
		}

		msg := tu.Message(tu.ID(update.Message.Chat.ID), text)
		_, sendErr := h.bot.SendMessage(ctx, msg)
		if sendErr != nil {
			log.Println("send error:", sendErr)
		}
		return
	}

	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		fmt.Sprintf("🐻 Медведь\nID: %d\nИмя: %s\n🍯 Мёд: %d", user.ID, user.UserName, user.Honey),
	)
	_, err = h.bot.SendMessage(ctx, msg)
	if err != nil {
		log.Println("send error:", err)
	}
}
