package bad_bees_tg_transport

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	bad_bees_service "github.com/toxagpark/HoneyGame/internal/feature/bad_bees/service"
)

type Handler struct {
	bot *telego.Bot
	// responseChatID — игровой чат из TG_CHAT_ID: сюда воркер пишет про налёт пчёл.
	responseChatID int64
	// interval — период налёта из BAD_BEES_INTERVAL.
	interval time.Duration
	service  service
	rand     *rand.Rand
}

type service interface {
	StingIdle(
		ctx context.Context,
	) ([]bad_bees_service.Sting, error)
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

// Run запускает воркер злых пчёл: раз в interval ищет АФК-медведей — без боёв
// и свежих вызовов за последнее время, отнимает у них процент мёда и пишет
// в игровой чат.
func (h *Handler) Run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.sting(ctx)
		}
	}
}

func (h *Handler) sting(ctx context.Context) {
	stings, err := h.service.StingIdle(ctx)
	if err != nil {
		log.Println("bad bees:", err)
		return
	}
	if len(stings) == 0 {
		log.Println("bad bees: все медведи дрались, пчёлы улетели ни с чем")
		return
	}

	msg := tu.Message(tu.ID(h.responseChatID), h.renderMessage(stings))
	if _, err := h.bot.SendMessage(ctx, msg); err != nil {
		log.Println("bad bees send error:", err)
	}
}

// renderMessage собирает сообщение: случайная шапка и список укушенных.
func (h *Handler) renderMessage(stings []bad_bees_service.Sting) string {
	var sb strings.Builder
	sb.WriteString(badBeesHeader(h.rand.Intn(len(badBeesHeaders))))
	sb.WriteString("\n\nКого покусали:\n")

	for _, s := range stings {
		sb.WriteString(fmt.Sprintf(
			"🐻 %s — 🍯 −%d (всего 🍯 %d)\n",
			s.UserName, s.Honey, s.NewHoney,
		))
	}

	return strings.TrimRight(sb.String(), "\n")
}

// badBeesHeader возвращает случайную шапку налёта по индексу.
func badBeesHeader(i int) string {
	return badBeesHeaders[i]
}

// badBeesHeaders — прикольные шапки злых пчёл, выбирается случайная.
var badBeesHeaders = []string{
	"🐝 Прилетели злые пчёлы и покусали тех, кто не дрался!",
	"🐝 Пчёлы провели рейд по лежебокам — счёт выставлен!",
	"🚨 Пасека объявила охоту на АФК-медведей!",
	"🐝 Не дрался — плати! Пчёлы забрали свою долю мёда.",
	"🔥 Злые пчёлы навестили всех, кто отсиживался в берлоге!",
}
