package users_tg_transport

import (
	"context"
	"fmt"
	"log"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (h *Handler) HandleCreateOrUpdateUser(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	newUser := domain.NewUserUninitialized(update.Message.From.ID, update.Message.From.Username)

	user, err := h.service.CreateOrUpdateUser(ctx, newUser)
	if err != nil {
		msg := tu.Message(tu.ID(update.Message.Chat.ID), "Ошибка сохранения пользователя(")
		log.Println(err)
		_, sendErr := h.bot.SendMessage(ctx, msg)
		if sendErr != nil {
			log.Println("send error:", sendErr)
		}
		return
	}

	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		fmt.Sprintf("🐻 Ты медведь!\nID: %d\nИмя: %s\n🍯 Мёд: %d", user.ID, user.UserName, user.Honey),
	)
	_, err = h.bot.SendMessage(ctx, msg)
	if err != nil {
		log.Println("send error:", err)
	}
}
