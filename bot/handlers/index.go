package handlers

import (
	"log"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers/ads"
	"tgbotBARAHOLKA/bot/handlers/profile"
	"tgbotBARAHOLKA/bot/handlers/start"
	"tgbotBARAHOLKA/bot/handlers/verification"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var nameHandlers = map[string]func(*tgbotapi.Update, *context.Context, int64){
	"start":        start.Handle,
	"ads":          ads.Handle,
	"profile":      profile.Handle,
	"Verification": verification.Handle,
}

func HandleUpdate(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64

	if update.Message != nil {
		userID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
	} else {
		log.Println("Ошибка: обновление не содержит данных о пользователе")
		return
	}

	state := context.GetUserState(userID, ctx)
	if update.Message != nil {
		switch update.Message.Command() {
		case "start":
			state.MessageID = 0
			context.ClearAllUserData(userID, ctx)
			start.HandleStartCommand(update, ctx)
			return
		}
	}
	if state.Level != 0 {
		if handler, exists := nameHandlers[state.Name]; exists {
			handler(update, ctx, userID)
		} else {
			start.Handle(update, ctx, userID)
		}
	} else {
		if update.CallbackQuery != nil {
			userId := update.CallbackQuery.From.ID
			switch strings.Split(update.CallbackQuery.Data, "_")[0] {
			case "adsMenu":
				context.UpdateUserName(userId, ctx, "ads")
				ads.HandleMenu(update, ctx)
			case "StartMenu":
				context.UpdateUserName(userId, ctx, "start")
				start.HandleStartCommand(update, ctx)
			case "AddAds":
				ads.HandleSelectADS(update, ctx)
			case "AdsHistory":
				ads.HandleSelectADSHistory(update, ctx)
			case "profile":
				context.UpdateUserName(userId, ctx, "profile")
				profile.HandleProfile(update, ctx)
			case "+balance":
				profile.HandleSelectPaymentMetod(update, ctx)
			case "Docs":
				start.HandleSelectDocs(update, ctx)
			case "Transfer":
				profile.HandleDoPayment(update, ctx)
			case "Verification":
				if len(strings.Split(update.CallbackQuery.Data, "_")) == 2 && strings.Split(update.CallbackQuery.Data, "_")[1] == strconv.Itoa(state.MessageID) {
					context.UpdateUserName(userId, ctx, "Verification")
					verification.HandleVerification(update, ctx)
				}
			}

		}
	}

}
