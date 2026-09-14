package fight_tg_transport

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

func (h *Handler) reply(ctx context.Context, chatID int64, text string) {
	h.send(ctx, tu.Message(tu.ID(chatID), text))
}

func (h *Handler) send(ctx context.Context, msg *telego.SendMessageParams) {
	if _, err := h.bot.SendMessage(ctx, msg); err != nil {
		log.Println("send error:", err)
	}
}

func (h *Handler) answerCallback(ctx context.Context, queryID string, text string) {
	if err := h.bot.AnswerCallbackQuery(ctx, tu.CallbackQuery(queryID).WithText(text)); err != nil {
		log.Println("answer callback error:", err)
	}
}

func (h *Handler) editMessage(ctx context.Context, query telego.CallbackQuery, text string) {
	msg := query.Message
	if msg == nil {
		return
	}
	if _, err := h.bot.EditMessageText(ctx, tu.EditMessageText(msg.GetChat().ChatID(), msg.GetMessageID(), text)); err != nil {
		log.Println("edit message error:", err)
	}
}

func parseAmount(text string) (int64, error) {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) < 2 {
		return 0, errors.New("amount is required")
	}
	return strconv.ParseInt(parts[1], 10, 64)
}

func parseCallbackID(data string, prefix string) (int, bool) {
	value, ok := strings.CutPrefix(data, prefix)
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return id, true
}

func tryStartText(err error) string {
	switch {
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrUserHoneyNotFound):
		return "Ты ещё не зарегистрирован 🐻 Напиши /start"
	default:
		log.Println(err)
		return "Ошибка( Попробуйте прописать /start"
	}
}
