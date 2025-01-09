package switchState

import (
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(update *tgbotapi.Update, ctx *context.Context) {
	handlers.HandleUpdate(update, ctx)
}
