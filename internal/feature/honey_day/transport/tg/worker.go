package honey_day_tg_transport

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	honey_day_service "github.com/toxagpark/HoneyGame/internal/feature/honey_day/service"
)

type Handler struct {
	bot *telego.Bot
	// responseChatID — игровой чат из TG_CHAT_ID: сюда воркер пишет про медовый день.
	responseChatID int64
	// interval — период медового дня из HONEY_DAY_INTERVAL.
	interval time.Duration
	service  service
	rand     *rand.Rand
}

type service interface {
	GiveHoneyDay(
		ctx context.Context,
	) ([]honey_day_service.HoneyDayGift, error)
}

func NewHandler(
	bot *telego.Bot,
	responseChatID int64,
	interval time.Duration,
	service service,
) *Handler {
	return &Handler{
		bot:            bot,
		responseChatID: responseChatID,
		interval:       interval,
		service:        service,
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Run запускает воркер медового дня: раз в interval раздаёт каждому свой
// подарок и пишет в игровой чат сообщение с итогами.
func (h *Handler) Run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.give(ctx)
		}
	}
}

func (h *Handler) give(ctx context.Context) {
	gifts, err := h.service.GiveHoneyDay(ctx)
	if err != nil {
		log.Println("honey day:", err)
		return
	}
	if len(gifts) == 0 {
		log.Println("honey day: no players yet, skipping")
		return
	}

	msg := tu.Message(tu.ID(h.responseChatID), h.renderMessage(gifts))
	if _, err := h.bot.SendMessage(ctx, msg); err != nil {
		log.Println("honey day send error:", err)
	}
}

// renderMessage собирает сообщение: случайная шапка, разыгранный улей и список.
func (h *Handler) renderMessage(gifts []honey_day_service.HoneyDayGift) string {
	total := int64(0)
	var sb strings.Builder
	sb.WriteString(honeyDayHeader(h.rand.Intn(len(honeyDayHeaders))))
	sb.WriteString("\n\n🍯 Медовый день: разыгран улей на 🍯 ")
	for _, g := range gifts {
		total += g.Honey
	}
	sb.WriteString(fmt.Sprintf("%d мёда!\n\nКому сколько:\n", total))

	for _, g := range gifts {
		sb.WriteString(fmt.Sprintf("🐻 %s — +%d мёда (всего 🍯 %d)\n", g.UserName, g.Honey, g.NewHoney))
	}

	return strings.TrimRight(sb.String(), "\n")
}

// honeyDayHeader возвращает случайную шапку медового дня по индексу.
func honeyDayHeader(i int) string {
	return honeyDayHeaders[i]
}

// honeyDayHeaders — прикольные шапки медового дня, выбирается случайная.
var honeyDayHeaders = []string{
	"🐝 Укуси меня пчела — медовый день пришёл!",
	"🐝 Пчёлы вернулись со смены, и у них всё с собой!",
	"🐻 Внимание! Улей открыл краны. Не зевай, медведь!",
	"🍯 Свежая партия мёда с пасеки уже здесь!",
	"🐝 Медовый день: пчёлы щедрые, лапы готовьте!",
}
