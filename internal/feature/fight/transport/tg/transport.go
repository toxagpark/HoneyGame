package fight_tg_transport

import (
	"context"
	"math/rand"
	"time"

	telego "github.com/mymmrac/telego"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	fight_service "github.com/toxagpark/HoneyGame/internal/feature/fight/service"
)

type Handler struct {
	bot *telego.Bot
	// responseChatID — игровой чат из TG_CHAT_ID: сюда бот отвечает всем игрокам.
	responseChatID int64
	service        service
	rand           *rand.Rand
}

// service — интерфейс фичи fight с точки зрения транспорта.
type service interface {
	CreateChallenge(
		ctx context.Context,
		tgUserID int64,
		amount int64,
	) (domain.ActiveChallenge, error)

	GetChallenges(
		ctx context.Context,
		tgUserID int64,
	) ([]domain.ActiveChallenge, error)

	CancelChallenge(
		ctx context.Context,
		tgUserID int64,
		challengeID int,
	) (domain.ActiveChallenge, error)

	ValidateAccept(
		ctx context.Context,
		tgUserID int64,
		challengeID int,
	) error

	AcceptChallenge(
		ctx context.Context,
		tgUserID int64,
		challengeID int,
	) (fight_service.AcceptChallengeResult, error)
}

func NewHandler(
	bot *telego.Bot,
	responseChatID int64,
	service service,
) *Handler {
	return &Handler{
		bot:            bot,
		responseChatID: responseChatID,
		service:        service,
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
