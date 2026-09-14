package menu_tg_transport

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	"github.com/toxagpark/HoneyGame/internal/core/tgutil"
	menu_service "github.com/toxagpark/HoneyGame/internal/feature/menu/service"
)

const menuText = "🍯 Добро пожаловать в Медовое Логово!\n\nВыбирай, что делаешь, медведь:"

// HandleMenu открывает главное меню по команде /menu.
func (h *Handler) HandleMenu(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	msg := tu.Message(tu.ID(update.Message.Chat.ID), menuText)
	msg.WithReplyMarkup(menuKeyboard())
	tgutil.Send(ctx, h.bot, msg)
}

// HandleMenuCallback — единая точка входа для всех кнопок меню.
func (h *Handler) HandleMenuCallback(ctx context.Context, query telego.CallbackQuery) {
	switch {
	case query.Data == callbackMenuProfile:
		h.showProfile(ctx, query)
	case query.Data == callbackMenuFight:
		h.showFightBets(ctx, query)
	case query.Data == callbackMenuFights:
		h.showFights(ctx, query)
	case query.Data == callbackMenuMyCall:
		h.showMyChallenge(ctx, query)
	case query.Data == callbackMenuBack:
		h.showMenu(ctx, query)
	default:
		if pct, ok := parseBetPercent(query.Data); ok {
			h.acceptBet(ctx, query, pct)
			return
		}
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Странная кнопка(")
	}
}

func (h *Handler) showMenu(ctx context.Context, query telego.CallbackQuery) {
	tgutil.EditCallbackMessage(ctx, h.bot, query, menuText)
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, menuRows())
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Меню")
}

func (h *Handler) showProfile(ctx context.Context, query telego.CallbackQuery) {
	user, err := h.service.Profile(ctx, query.From.ID)
	if err != nil {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, tryStartText(err))
		return
	}

	text := fmt.Sprintf(
		"👤 Профиль медведя\n\n🆔 ID: %d\n🐻 Имя: %s\n🍯 Мёд: %d",
		user.ID, user.UserName, user.Honey,
	)
	tgutil.EditCallbackMessage(ctx, h.bot, query, text)
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Профиль")
}

// showFightBets предлагает выбор ставки: процент от баланса, округлённый вниз.
func (h *Handler) showFightBets(ctx context.Context, query telego.CallbackQuery) {
	user, err := h.service.Profile(ctx, query.From.ID)
	if err != nil {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, tryStartText(err))
		return
	}

	options := h.service.BetOptions(ctx, user.Honey)
	if len(options) == 0 {
		tgutil.EditCallbackMessage(ctx, h.bot, query,
			fmt.Sprintf("🍯 У тебя всего %d мёда — на бой не накопить.\n\nИди подрабатывай медведем 🐻‍❄️", user.Honey))
		tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Мало мёда")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("⚔️ Создание боя\n\n🍯 Твой баланс: %d\n\nВыбери ставку (%% от баланса, округляется вниз):\n", user.Honey))
	for _, opt := range options {
		sb.WriteString(fmt.Sprintf("• %d%% → 🍯 %d\n", opt.Percent, opt.Amount))
	}

	tgutil.EditCallbackMessage(ctx, h.bot, query, sb.String())
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, betKeyboard(options))
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Выбирай ставку")
}

func (h *Handler) showFights(ctx context.Context, query telego.CallbackQuery) {
	challenges, err := h.service.Fights(ctx, query.From.ID)
	if err != nil {
		tgutil.AnswerCallback(ctx, h.bot, query.ID, tryStartText(err))
		return
	}

	if len(challenges) == 0 {
		tgutil.EditCallbackMessage(ctx, h.bot, query, "🔍 Открытых вызовов нет 🐻\n\nСоздай свой через ⚔️ «Создать бой»!")
		tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Пусто")
		return
	}

	var sb strings.Builder
	sb.WriteString("🔍 Открытые вызовы:\n\n")
	buttons := [][]telego.InlineKeyboardButton{}
	for _, c := range challenges {
		sb.WriteString(fmt.Sprintf("№%d — 🐻 %s ставит 🍯 %d\n", c.ID, c.CreatorName, c.Amount))
		buttons = append(buttons, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprintf("🗡 Взять №%d · 🍯 %d", c.ID, c.Amount)).
				WithCallbackData(fmt.Sprintf("%s%d", fightAcceptPrefix, c.ID)),
		))
	}

	tgutil.EditCallbackMessage(ctx, h.bot, query, sb.String())
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, append(buttons, backRow()...))
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Вызовы")
}

// showMyChallenge показывает активный вызов игрока с кнопкой отзыва;
// если вызова нет — предлагает создать.
func (h *Handler) showMyChallenge(ctx context.Context, query telego.CallbackQuery) {
	challenge, err := h.service.MyChallenge(ctx, query.From.ID)
	if err != nil {
		if errors.Is(err, domain.ErrChallengeNotFound) {
			tgutil.EditCallbackMessage(ctx, h.bot, query,
				"📌 Активного вызова нет 🐻\n\nСоздай через «⚔️ Создать бой» — или командой /fight <ставка>")
			tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
			tgutil.AnswerCallback(ctx, h.bot, query.ID, "Вызова нет")
			return
		}
		tgutil.AnswerCallback(ctx, h.bot, query.ID, tryStartText(err))
		return
	}

	tgutil.EditCallbackMessage(ctx, h.bot, query,
		fmt.Sprintf("📌 Твой вызов №%d\n🍯 Ставка: %d\n⏳ Ждём соперника!", challenge.ID, challenge.Amount))
	// Отзыв идёт через общий fight-флоу: тот же callback, что и под сообщением вызова.
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, append([][]telego.InlineKeyboardButton{
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("✖️ Отозвать").WithCallbackData(fmt.Sprintf("fight:cancel:%d", challenge.ID)),
		),
	}, backRow()...))
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Твой вызов")
}

func (h *Handler) acceptBet(ctx context.Context, query telego.CallbackQuery, pct int) {
	challenge, err := h.service.CreateChallengeFromMenu(ctx, query.From.ID, pct)
	if err != nil {
		tgutil.EditCallbackMessage(ctx, h.bot, query, createChallengeErrorText(err))
		tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
		tgutil.AnswerCallback(ctx, h.bot, query.ID, "Не вышло")
		return
	}

	tgutil.EditCallbackMessage(ctx, h.bot, query,
		fmt.Sprintf("⚔️ Вызов №%d создан!\n🐻 Ставка: 🍯 %d\n⏳ Ждём соперника!", challenge.ID, challenge.Amount))
	tgutil.EditCallbackKeyboard(ctx, h.bot, query, backRow())
	tgutil.AnswerCallback(ctx, h.bot, query.ID, "Вызов создан! 🎉")
}

// --- кнопки и парсинг ---

func menuRows() [][]telego.InlineKeyboardButton {
	return [][]telego.InlineKeyboardButton{
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("👤 Профиль").WithCallbackData(callbackMenuProfile)),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("⚔️ Создать бой").WithCallbackData(callbackMenuFight),
			tu.InlineKeyboardButton("📌 Мой вызов").WithCallbackData(callbackMenuMyCall),
		),
		tu.InlineKeyboardRow(tu.InlineKeyboardButton("🔍 Смотреть бои").WithCallbackData(callbackMenuFights)),
	}
}

// menuKeyboard — готовая разметка главного меню.
func menuKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(menuRows()...)
}

func betKeyboard(options []menu_service.BetOption) [][]telego.InlineKeyboardButton {
	rows := [][]telego.InlineKeyboardButton{}
	for _, opt := range options {
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprintf("🍯 %d%% — %d мёда", opt.Percent, opt.Amount)).
				WithCallbackData(fmt.Sprintf("%s%d", callbackMenuBetFmt, opt.Percent)),
		))
	}
	return append(rows, backRow()...)
}

func parseBetPercent(data string) (int, bool) {
	value, ok := strings.CutPrefix(data, callbackMenuBetFmt)
	if !ok {
		return 0, false
	}
	pct, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return pct, true
}
