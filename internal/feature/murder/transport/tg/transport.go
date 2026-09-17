package murder_tg_transport

import (
	"context"

	telego "github.com/mymmrac/telego"
	murder_service "github.com/toxagpark/HoneyGame/internal/feature/murder/service"
)

type Handler struct {
	bot *telego.Bot
	// responseChatID — игровой чат из TG_CHAT_ID: сюда бот отвечает всем игрокам.
	responseChatID int64
	service        service
}

// service — интерфейс фичи murder с точки зрения транспорта.
type service interface {
	Rob(
		ctx context.Context,
		tgUserID int64,
		hives int,
		amount int64,
	) (murder_service.RobResult, error)
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
