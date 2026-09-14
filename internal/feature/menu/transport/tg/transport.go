package menu_tg_transport

import (
	telego "github.com/mymmrac/telego"
	menu_service "github.com/toxagpark/HoneyGame/internal/feature/menu/service"
)

type Handler struct {
	bot     *telego.Bot
	service *menu_service.Service
}

func NewHandler(
	bot *telego.Bot,
	service *menu_service.Service,
) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}
