// Package tgutil — общие помощники Telegram-транспорта: отправка, редактирование,
// ответы на callback, парсинг callback-data. Не содержит бизнес-логики.
package tgutil

import (
	"context"
	"log"
	"strconv"
	"strings"

	telego "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Send(ctx context.Context, bot *telego.Bot, msg *telego.SendMessageParams) {
	if _, err := bot.SendMessage(ctx, msg); err != nil {
		log.Println("send error:", err)
	}
}

func Reply(ctx context.Context, bot *telego.Bot, chatID int64, text string) {
	Send(ctx, bot, tu.Message(tu.ID(chatID), text))
}

func AnswerCallback(ctx context.Context, bot *telego.Bot, queryID string, text string) {
	if err := bot.AnswerCallbackQuery(ctx, tu.CallbackQuery(queryID).WithText(text)); err != nil {
		log.Println("answer callback error:", err)
	}
}

// EditCallbackMessage редактирует сообщение, к которому была прижата кнопка.
func EditCallbackMessage(ctx context.Context, bot *telego.Bot, query telego.CallbackQuery, text string) {
	msg := query.Message
	if msg == nil {
		return
	}
	if _, err := bot.EditMessageText(ctx, tu.EditMessageText(msg.GetChat().ChatID(), msg.GetMessageID(), text)); err != nil {
		log.Println("edit message error:", err)
	}
}

// EditCallbackKeyboard меняет inline-клавиатуру у сообщения с кнопкой.
func EditCallbackKeyboard(ctx context.Context, bot *telego.Bot, query telego.CallbackQuery, rows [][]telego.InlineKeyboardButton) {
	msg := query.Message
	if msg == nil {
		return
	}
	params := tu.EditMessageReplyMarkup(msg.GetChat().ChatID(), msg.GetMessageID(), tu.InlineKeyboard(rows...))
	if _, err := bot.EditMessageReplyMarkup(ctx, params); err != nil {
		log.Println("edit reply markup error:", err)
	}
}

func ParseCallbackID(data string, prefix string) (int, bool) {
	value, ok := strings.CutPrefix(data, prefix)
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return id, true
}
