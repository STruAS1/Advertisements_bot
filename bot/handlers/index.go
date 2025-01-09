package handlers

import (
	"log"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers/ads"
	"tgbotBARAHOLKA/bot/handlers/start"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var nameHandlers = map[string]func(*tgbotapi.Update, *context.Context, int64){
	"start": start.Handle,
	"ads":   ads.Handle,
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
	if state.Level != 0 {
		if handler, exists := nameHandlers[state.Name]; exists {
			handler(update, ctx, userID)
		} else {
			start.Handle(update, ctx, userID)
		}
	} else {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				start.HandleStartCommand(update, ctx)
			}
		}
		if update.CallbackQuery != nil {
			userId := update.CallbackQuery.From.ID
			switch update.CallbackQuery.Data {
			case "adsMenu":
				context.UpdateUserName(userId, ctx, "ads")
				ads.HandleMenu(update, ctx)
			case "StartMenu":
				start.HandleStartCommand(update, ctx)
			case "AddAds":
				ads.HandleSelectADS(update, ctx)
			}

		}
	}

}
