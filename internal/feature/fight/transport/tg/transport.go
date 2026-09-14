package fight_tg_transport

import (
	"context"

	telego "github.com/mymmrac/telego"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
	fight_service "github.com/toxagpark/HoneyGame/internal/feature/fight/service"
)

type Handler struct {
	bot     *telego.Bot
	service service
}

// service — интерфейс фичи fight с точки зрения транспорта.
type service interface {
	CreateChallenge(
		ctx context.Context,
		tgChatID int64,
		amount int64,
	) (domain.ActiveChallenge, error)

	GetChallenges(
		ctx context.Context,
		tgChatID int64,
	) ([]domain.ActiveChallenge, error)

	CancelChallenge(
		ctx context.Context,
		tgChatID int64,
		challengeID int,
	) (domain.ActiveChallenge, error)

	ValidateAccept(
		ctx context.Context,
		tgChatID int64,
		challengeID int,
	) error

	AcceptChallenge(
		ctx context.Context,
		tgChatID int64,
		challengeID int,
	) (fight_service.AcceptChallengeResult, error)
}

func NewHandler(
	bot *telego.Bot,
	service service,
) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}
