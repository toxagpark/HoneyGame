package fight_tg_transport

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	tgutil "github.com/toxagpark/HoneyGame/internal/core/transport/tg/util"
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
		tgutil.Reply(ctx, h.bot, h.responseChatID, "Укажи ставку числом: /fight 10 🐻")
		return
	}

	challenge, err := h.service.CreateChallenge(ctx, update.Message.From.ID, amount)
	if err != nil {
		tgutil.Reply(ctx, h.bot, h.responseChatID, createChallengeErrorText(err))
		return
	}

	msg := tu.Message(
		tu.ID(h.responseChatID),
		fmt.Sprintf("⚔️ Вызов №%d создан!\n🐻 %s ставит 🍯 %d\nКто смелый?!", challenge.ID, challenge.CreatorName, challenge.Amount),
	)
	// Свои вызовы принимает только соперник — из меню и /fights, поэтому кнопка «Взять» тут не нужна.
	msg.WithReplyMarkup(
		tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("✖️ Отозвать").WithCallbackData(fmt.Sprintf("%s%d", callbackPrefixCancel, challenge.ID)),
			),
		),
	)
	tgutil.Send(ctx, h.bot, msg)
}

func (h *Handler) HandleListChallenges(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	challenges, err := h.service.GetChallenges(ctx, update.Message.From.ID)
	if err != nil {
		tgutil.Reply(ctx, h.bot, h.responseChatID, tryStartText(err))
		return
	}

	renderChallenges(ctx, h, h.responseChatID, challenges)
}

// renderChallenges рисует список вызовов: текст + кнопки «взять».
func renderChallenges(
	ctx context.Context,
	h *Handler,
	chatID int64,
	challenges []domain.ActiveChallenge,
) {
	if len(challenges) == 0 {
		tgutil.Reply(ctx, h.bot, chatID, "Открытых вызовов нет 🐻 Создай свой: /fight <ставка>")
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

	msg := tu.Message(tu.ID(chatID), sb.String())
	msg.WithReplyMarkup(tu.InlineKeyboard(buttons...))
	tgutil.Send(ctx, h.bot, msg)
}

func (h *Handler) HandleCancelCallback(ctx context.Context, query telego.CallbackQuery) {
	challengeID, ok := tgutil.ParseCallbackID(query.Data, callbackPrefixCancel)
	if !ok {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Странная кнопка(")
		return
	}

	challenge, err := h.service.CancelChallenge(ctx, query.From.ID, challengeID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			tgutil.AnswerCallback(ctx, h.bot, query.ID, "Вызов уже не актуален 🐻")
			return
		}
		if errors.Is(err, domain.ErrChallengeNotOwner) {
			tgutil.AnswerCallback(ctx, h.bot, query.ID, "Отозвать может только создатель вызова 🐻")
			return
		}
		tgutil.EditCallbackMessage(ctx, h.bot, query, tryStartText(err))
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Ошибка отзыва")
		return
	}

	tgutil.EditCallbackMessage(ctx, h.bot, query, fmt.Sprintf("✖️ Вызов №%d на 🍯 %d отозван 🐻", challenge.ID, challenge.Amount))
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Вызов отозван")
}

func (h *Handler) HandleAcceptCallback(ctx context.Context, query telego.CallbackQuery) {
	challengeID, ok := tgutil.ParseCallbackID(query.Data, callbackPrefixAccept)
	if !ok {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Странная кнопка(")
		return
	}

	// Ошибки отсекаем до анимации, чтобы не «играть» заведомый провал.
	if err := h.service.ValidateAccept(ctx, query.From.ID, challengeID); err != nil {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, validateAcceptErrorText(err))
		return
	}

	// Анимация «идёт игра» — редактируем сообщение с кнопками.
	playFightAnimation(ctx, h, query)

	result, err := h.service.AcceptChallenge(ctx, query.From.ID, challengeID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrChallengeNotFound):
			tgutil.EditCallbackMessage(ctx, h.bot, query, "⌛️ Вызов уже забрали другим медведем 🐻")
		case errors.Is(err, domain.ErrNotEnoughHoney):
			tgutil.EditCallbackMessage(ctx, h.bot, query, "🤷 Не хватило мёда на бой. Попробуйте прописать /start")
		default:
			tgutil.EditCallbackMessage(ctx, h.bot, query, "Ошибка боя( Попробуйте прописать /start")
		}
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Бой не состоялся")
		return
	}

	announceResult(ctx, h, query, challengeID, result)
}

func playFightAnimation(ctx context.Context, h *Handler, query telego.CallbackQuery) {
	frames := []string{
		"🐻‍❄️ Мишки сходятся...",
		"💥 Удар!",
		"🍯 Мёд летит во все стороны!",
		"🤜💫🤛",
	}

	for _, frame := range frames {
		tgutil.EditCallbackMessage(ctx, h.bot, query, frame)
		time.Sleep(800 * time.Millisecond)
	}
}

func announceResult(ctx context.Context, h *Handler, query telego.CallbackQuery, challengeID int, result fight_service.AcceptChallengeResult) {
	// Исход боя адресно знает только тот, кто жал кнопку; остальной чат видит нейтральное сообщение.
	if result.AcceptorWon {
		tgutil.EditCallbackMessage(ctx, h.bot, query, "🏆 Ты победил!\n🍯 +"+fmt.Sprint(result.Amount)+" мёда")
	} else {
		tgutil.EditCallbackMessage(ctx, h.bot, query, "💀 Ты проиграл...\n🍯 -"+fmt.Sprint(result.Amount)+" мёда")
	}
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Итог боя")

	tgutil.Reply(ctx, h.bot, h.responseChatID, fmt.Sprintf(
		"⚔️ Вызов №%d на 🍯 %d разыгран\nМёд нашёл своего медведя 🐻",
		challengeID, result.Amount,
	))
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
	case errors.Is(err, domain.ErrNotEnoughHoney):
		return "🍯 У тебя не хватает мёда на этот вызов"
	default:
		return tryStartText(err)
	}
}
