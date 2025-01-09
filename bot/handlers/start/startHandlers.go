package start

import (
	"tgbotBARAHOLKA/bot/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	switch state.Level {
	case 1:
		if update.Message != nil {
			HandlePhoneNumberRequest(update, ctx)
		}
	case 2:
		handleLvl2(update, ctx)
	}
}

func handleLvl2(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "cehk_sub":
			HandleSubscriptionCheck(update, ctx)
		}
	}
}
