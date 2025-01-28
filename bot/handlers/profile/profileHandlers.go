package profile

import (
	"strings"
	"tgbotBARAHOLKA/bot/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	switch state.Level {
	case 1:
		handleLvl1(update, ctx, userID)
	case 2:
		handleLvl2(update, ctx, userID)
	case 3:
		handleLvl3(update, ctx, userID)
	case 4:
		handleLvl4(update, ctx, userID)
		// case 3:
		// 	handleLvl3(update, ctx, userID)
		// case 4:
		// 	handleLvl4(update, ctx, userID)
		// case 5:
		// 	handelLvl5(update, ctx, userID)
		// case 6:
		// 	handelLvl6(update, ctx)
		// case 7:
		// 	handelLvl7(update, ctx, userID)
		// case 8:
		// 	handelLvl8(update, ctx, userID)
		// case 9:
		// 	handelLvl9(update, ctx, userID)
	}
}

func handleLvl1(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		println(update.CallbackQuery.Data)
		if len(data) == 1 {
			switch data[0] {
			case "back":
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "Payment")
				HandleProfile(update, ctx)
			}
		} else if len(data) == 2 && data[0] == "payment" {
			println(data[0])
			HandlePaymentEntryAmount(update, ctx)
		}
	}
}
func handleLvl2(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		if len(data) == 1 {
			switch data[0] {
			case "back":
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "Payment")
				HandleProfile(update, ctx)
			}
		}
	} else if update.Message != nil {
		HandleShowMetods(update, ctx)
	}
}
func handleLvl3(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		if len(data) == 1 {
			switch data[0] {
			case "back":
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "Payment")
				HandleProfile(update, ctx)
			case "confirm":
				HandeleConfirmPayment(update, ctx)
			}
		}
	}
}

func handleLvl4(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			state := context.GetUserState(userID, ctx)
			delete(state.Data, "transfer")
			HandleProfile(update, ctx)
		default:
			HandleDoPayment(update, ctx)
		}
	}
	if update.Message != nil {
		HandleDoPayment(update, ctx)
	}
}
