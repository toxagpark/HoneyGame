package users_tg_transport

import (
	"context"
	"fmt"
	"log"
	"strings"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// HandleGetTopUsers выводит топ игроков: имя — мёд, от большего к меньшему.
func (h *Handler) HandleGetTopUsers(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	users, err := h.service.GetTopUsers(ctx)
	if err != nil {
		log.Println(err)
		msg := tu.Message(tu.ID(h.responseChatID), "Ошибка получения топа(")
		if _, sendErr := h.bot.SendMessage(ctx, msg); sendErr != nil {
			log.Println("send error:", sendErr)
		}
		return
	}

	var sb strings.Builder
	sb.WriteString("🍯 Топ медведей по мёду:\n")
	for i, u := range users {
		sb.WriteString(fmt.Sprintf("%d. %s — 🍯 %d\n", i+1, u.UserName, u.Honey))
	}

	msg := tu.Message(tu.ID(h.responseChatID), sb.String())
	if _, sendErr := h.bot.SendMessage(ctx, msg); sendErr != nil {
		log.Println("send error:", sendErr)
	}
}
