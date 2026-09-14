package users_tg_transport

import (
	"context"

	telego "github.com/mymmrac/telego"
	"github.com/toxagpark/HoneyGame/internal/core/domain"
)

type Handler struct {
	bot     *telego.Bot
	service service
}

type service interface {
	CreateOrUpdateUser(
		ctx context.Context,
		newUser domain.User,
	) (domain.User, error)

	GetUser(
		ctx context.Context,
		tgChatID int64,
	) (domain.User, error)
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
