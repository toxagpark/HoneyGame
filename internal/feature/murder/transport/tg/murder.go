package murder_tg_transport

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	tgutil "github.com/toxagpark/HoneyGame/internal/core/transport/tg/util"
	murder_service "github.com/toxagpark/HoneyGame/internal/feature/murder/service"
)

// animationFrameDelay — задержка между кадрами анимации ограбления.
const animationFrameDelay = 800 * time.Millisecond

func (h *Handler) HandleMurder(ctx context.Context, update *telego.Update) {
	if update.Message == nil {
		return
	}

	hives, amount, err := parseArgs(update.Message.Text)
	if err != nil {
		tgutil.Reply(ctx, h.bot, h.responseChatID, rulesText())
		return
	}

	result, err := h.service.Rob(ctx, update.Message.From.ID, hives, amount)
	if err != nil {
		tgutil.Reply(ctx, h.bot, h.responseChatID, robErrorText(err))
		return
	}

	playRobAnimation(ctx, h, hives, result)
	announceResult(ctx, h, hives, amount, result)
}

// parseArgs разбирает «/murder <ульи> <ставка>»; без аргументов — правила.
func parseArgs(text string) (int, int64, error) {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) < 3 {
		return 0, 0, errors.New("hives and stake are required")
	}

	hives, err := parsePositiveInt(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse hives: %w", err)
	}

	amount, err := parsePositiveInt(parts[2])
	if err != nil {
		return 0, 0, fmt.Errorf("parse stake: %w", err)
	}

	return hives, int64(amount), nil
}

// parsePositiveInt читает целое число строго из одного токена: «5abc» и «5.5» — не числа.
func parsePositiveInt(token string) (int, error) {
	n, err := strconv.Atoi(token)
	if err != nil || n <= 0 {
		return 0, errors.New("not a positive integer")
	}
	return n, nil
}

// rulesText — правила по /murder без аргументов.
func rulesText() string {
	minHives, maxHives := murder_service.HivesBounds()
	return fmt.Sprintf(
		"🐻 ОГРАБЛЕНИЕ УЛЬЯ\n\n"+
			"На поляне стоят ульи — мёд спрятан только в одном!\n"+
			"Ставишь мёд, мишка тычет лапу в случайный улей.\n\n"+
			"🍯 Угадал — забираешь ставку × кол-во ульев\n"+
			"🐝 Не угадал — пчёлы улетают с твоим мёдом\n\n"+
			"Играешь так:\n"+
			"/murder <ульи %d–%d> <ставка>\n"+
			"Например: /murder 5 10\n\n"+
			"Шанс — 1 к кол-ву ульев, зато и выигрыш жирный!",
		minHives, maxHives,
	)
}

// playRobAnimation показывает пасеку, ведёт лапу медведя к выбранному улью,
// продавливает улей и вскрывает итог: где мёд.
func playRobAnimation(
	ctx context.Context,
	h *Handler,
	hives int,
	result murder_service.RobResult,
) {
	frames := buildFrames(result.PawHive, result.HoneyHive, result.Honey > 0)

	messageID, err := tgutil.SendMessage(
		ctx,
		h.bot,
		tu.Message(tu.ID(h.responseChatID), sceneText(frames[0], hives)),
	)
	if err != nil {
		return
	}

	for _, frame := range frames[1:] {
		time.Sleep(animationFrameDelay)
		tgutil.EditMessage(ctx, h.bot, h.responseChatID, messageID, sceneText(frame, hives))
	}
}

// --- итог ---

// announceResult пишет итог ограбления одним сообщением.
func announceResult(
	ctx context.Context,
	h *Handler,
	hives int,
	amount int64,
	result murder_service.RobResult,
) {
	var text string
	if result.Honey > 0 {
		text = fmt.Sprintf(
			"🏆 @%s вскрыл улей %s — а там МЁД!\n🍯 %d × %d = +%d мёда (всего: %d)",
			result.UserName, keycap(result.PawHive), amount, hives, result.Honey, result.NewHoney,
		)
	} else {
		text = fmt.Sprintf(
			"💀 @%s ткнул лапу в %s, а мёд лежал в %s\n🐝 Ставка ушла пчёлам (всего: %d)",
			result.UserName, keycap(result.PawHive), keycap(result.HoneyHive), result.NewHoney,
		)
	}
	tgutil.Reply(ctx, h.bot, h.responseChatID, text)
}

// --- тексты ошибок ---

func robErrorText(err error) string {
	minHives, maxHives := murder_service.HivesBounds()
	switch {
	case errors.Is(err, domain.ErrWrongHives):
		return fmt.Sprintf("Ульев должно быть от %d до %d 🐻", minHives, maxHives)
	case errors.Is(err, domain.ErrWrongAmount):
		return "Ставка должна быть целым числом больше нуля 🐻"
	case errors.Is(err, domain.ErrNotEnoughHoney):
		return "🤷 Не хватает мёда на такую ставку 🐻 Пополни запасы медовым днём"
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrUserHoneyNotFound):
		return "Ты ещё не зарегистрирован 🐻 Напиши /start"
	default:
		return "Ошибка( Попробуйте прописать /start"
	}
}
