package users_tg_transport

import (
	"context"

	telego "github.com/mymmrac/telego"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type Handler struct {
	bot *telego.Bot
	// responseChatID — игровой чат из TG_CHAT_ID: сюда бот отвечает всем игрокам.
	responseChatID int64
	service        service
}

type service interface {
	CreateOrUpdateUser(
		ctx context.Context,
		newUser domain.User,
	) (domain.User, error)

	GetUser(
		ctx context.Context,
		tgUserID int64,
	) (domain.User, error)

	GetTopUsers(
		ctx context.Context,
	) ([]domain.User, error)
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
	}
}
