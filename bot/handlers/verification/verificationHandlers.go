package verification

import (
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers/start"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	switch state.Level {
	case 1:
		handleLvl1(update, ctx)
	}

}

func handleLvl1(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			start.HandleStartCommand(update, ctx)
		default:
			HandleVerification(update, ctx)
		}
	}
	if update.Message != nil {
		HandleVerification(update, ctx)
	}
}
