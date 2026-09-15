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

// HandleGetMe показывает профиль вызвавшего: внутренний ID, имя, баланс мёда.
func (h *Handler) HandleGetMe(ctx context.Context, update *telego.Update) {
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

		msg := tu.Message(tu.ID(h.responseChatID), text)
		if _, sendErr := h.bot.SendMessage(ctx, msg); sendErr != nil {
			log.Println("send error:", sendErr)
		}
		return
	}

	msg := tu.Message(
		tu.ID(h.responseChatID),
		fmt.Sprintf("🐻 Ты медведь!\nID: %d\nИмя: %s\n🍯 Мёд: %d", user.ID, user.UserName, user.Honey),
	)
	if _, sendErr := h.bot.SendMessage(ctx, msg); sendErr != nil {
		log.Println("send error:", sendErr)
	}
}
