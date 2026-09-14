package fight_tg_transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	fight_service "github.com/toxagpark/HoneyGame/internal/feature/fight/service"
)

// callback data: fight:accept:<challengeID>, fight:cancel:<challengeID>
const callbackPrefixAccept = "fight:accept:"
const callbackPrefixCancel = "fight:cancel:"

func (h *Handler) HandleCreateChallenge(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	amount, err := parseAmount(update.Message.Text)
	if err != nil {
		h.reply(ctx, update.Message.Chat.ID, "Укажи ставку числом: /fight 10 🐻")
		return
	}

	challenge, err := h.service.CreateChallenge(ctx, update.Message.From.ID, amount)
	if err != nil {
		h.reply(ctx, update.Message.Chat.ID, createChallengeErrorText(err))
		return
	}

	msg := tu.Message(
		tu.ID(update.Message.Chat.ID),
		fmt.Sprintf("⚔️ Вызов №%d создан!\n🐻 %s ставит 🍯 %d\nКто смелый?!", challenge.ID, challenge.CreatorName, challenge.Amount),
	)
	// Свои вызовы принимает только соперник — из /fights, поэтому кнопка «Взять» тут не нужна.
	msg.WithReplyMarkup(
		tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("✖️ Отозвать").WithCallbackData(fmt.Sprintf("%s%d", callbackPrefixCancel, challenge.ID)),
			),
		),
	)
	h.send(ctx, msg)
}

func (h *Handler) HandleListChallenges(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	challenges, err := h.service.GetChallenges(ctx, update.Message.From.ID)
	if err != nil {
		log.Println(err)
		h.reply(ctx, update.Message.Chat.ID, tryStartText(err))
		return
	}

	if len(challenges) == 0 {
		h.reply(ctx, update.Message.Chat.ID, "Открытых вызовов нет 🐻 Создай свой: /fight <ставка>")
		return
	}

	var sb strings.Builder
	sb.WriteString("⚔️ Открытые вызовы:\n\n")
	buttons := [][]telego.InlineKeyboardButton{}
	for _, c := range challenges {
		sb.WriteString(fmt.Sprintf("№%d — 🐻 %s ставит 🍯 %d\n", c.ID, c.CreatorName, c.Amount))
		buttons = append(buttons, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprintf("🗡 №%d · %s · 🍯 %d", c.ID, c.CreatorName, c.Amount)).
				WithCallbackData(fmt.Sprintf("%s%d", callbackPrefixAccept, c.ID)),
		))
	}

	msg := tu.Message(tu.ID(update.Message.Chat.ID), sb.String())
	msg.WithReplyMarkup(tu.InlineKeyboard(buttons...))
	h.send(ctx, msg)
}

func (h *Handler) HandleCancelCallback(ctx context.Context, query telego.CallbackQuery) {
	challengeID, ok := parseCallbackID(query.Data, callbackPrefixCancel)
	if !ok {
		h.answerCallback(ctx, query.ID, "Странная кнопка(")
		return
	}

	challenge, err := h.service.CancelChallenge(ctx, query.From.ID, challengeID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			h.answerCallback(ctx, query.ID, "Вызов уже не актуален 🐻")
			return
		}
		if errors.Is(err, domain.ErrChallengeNotOwner) {
			h.answerCallback(ctx, query.ID, "Отозвать может только создатель вызова 🐻")
			return
		}
		log.Println(err)
		h.editMessage(ctx, query, tryStartText(err))
		h.answerCallback(ctx, query.ID, "Ошибка отзыва")
		return
	}

	h.editMessage(ctx, query, fmt.Sprintf("✖️ Вызов №%d на 🍯 %d отозван 🐻", challenge.ID, challenge.Amount))
	h.answerCallback(ctx, query.ID, "Вызов отозван")
}

func (h *Handler) HandleAcceptCallback(ctx context.Context, query telego.CallbackQuery) {
	challengeID, ok := parseCallbackID(query.Data, callbackPrefixAccept)
	if !ok {
		h.answerCallback(ctx, query.ID, "Странная кнопка(")
		return
	}

	// Ошибки отсекаем до анимации, чтобы не «играть» заведомый провал.
	if err := h.service.ValidateAccept(ctx, query.From.ID, challengeID); err != nil {
		log.Println(err)
		h.answerCallback(ctx, query.ID, validateAcceptErrorText(err))
		return
	}

	// Анимация «идёт игра» — редактируем сообщение с кнопками.
	h.playFightAnimation(ctx, query)

	result, err := h.service.AcceptChallenge(ctx, query.From.ID, challengeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrChallengeNotFound):
			h.editMessage(ctx, query, "⌛️ Вызов уже забрали другим медведем 🐻")
		case errors.Is(err, domain.ErrNotEnoughHoney):
			h.editMessage(ctx, query, "🤷 Не хватило мёда на бой. Попробуйте прописать /start")
		default:
			log.Println(err)
			h.editMessage(ctx, query, "Ошибка боя( Попробуйте прописать /start")
		}
		h.answerCallback(ctx, query.ID, "Бой не состоялся")
		return
	}

	h.announceResult(ctx, query, result)
}

func (h *Handler) playFightAnimation(ctx context.Context, query telego.CallbackQuery) {
	frames := []string{
		"🐻‍❄️ Мишки сходятся...",
		"💥 Удар!",
		"🍯 Мёд летит во все стороны!",
		"🤜💫🤛",
	}

	for _, frame := range frames {
		h.editMessage(ctx, query, frame)
		time.Sleep(800 * time.Millisecond)
	}
}

func (h *Handler) announceResult(ctx context.Context, query telego.CallbackQuery, result fight_service.AcceptChallengeResult) {
	var acceptorText string
	if result.AcceptorWon {
		acceptorText = fmt.Sprintf("🏆 Ты победил!\n🍯 +%d мёда", result.Amount)
	} else {
		acceptorText = fmt.Sprintf("💀 Ты проиграл...\n🍯 -%d мёда", result.Amount)
	}
	h.editMessage(ctx, query, acceptorText)
	h.answerCallback(ctx, query.ID, "Итог боя")

	// Итог второму участнику отдельным сообщением.
	var opponentText string
	if result.AcceptorWon {
		opponentText = fmt.Sprintf("💀 Твой вызов приняли и ты проиграл...\n🍯 -%d мёда", result.Amount)
	} else {
		opponentText = fmt.Sprintf("🏆 Твой вызов приняли и ты победил!\n🍯 +%d мёда", result.Amount)
	}
	h.reply(ctx, result.OpponentTgChatID, opponentText)
}

// --- тексты ошибок ---

func createChallengeErrorText(err error) string {
	switch {
	case errors.Is(err, domain.ErrChallengeAlreadyExists):
		return "У тебя уже есть активный вызов! Отзови его или дождись боя 🐻"
	case errors.Is(err, domain.ErrWrongAmount):
		return "Ставка должна быть больше нуля 🐻"
	case errors.Is(err, domain.ErrNotEnoughHoney):
		return "Не хватает мёда на такую ставку 🍯🐻"
	default:
		return tryStartText(err)
	}
}

func validateAcceptErrorText(err error) string {
	switch {
	case errors.Is(err, domain.ErrChallengeNotFound):
		return "Вызов уже не актуален 🐻"
	case errors.Is(err, domain.ErrSelfChallenge):
		return "На себя драться нельзя 🐻"
	default:
		return tryStartText(err)
	}
}
