package ads

import (
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	switch state.Level {
	case 1:
		handleLvl1(update, ctx)
	case 2:
		handleLvl2(update, ctx, userID)
	case 3:
		handleLvl3(update, ctx, userID)
	case 4:
		handleLvl4(update, ctx, userID)
	case 5:
		handelLvl5(update, ctx, userID)
	case 6:
		handelLvl6(update, ctx)
	}
}

func handleLvl1(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		if len(data) == 1 {
			switch data[0] {
			case "back":
				HandleMenu(update, ctx)
			}
		} else if len(data) == 2 && data[0] == "newAds" {
			HandleAddAds(update, ctx, data[1])
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
				delete(state.Data, "AdsInputs")
				HandleSelectADS(update, ctx)
			case "AddPhoto":
				HandleAddPhoto(update, ctx)
			case "preViwe":
				HandlePreWive(update, ctx)
			case "Save":
				HandleSaveAds(update, ctx)
			}

		} else if len(data) == 2 && data[0] == "AddInput" {
			HandleAddInput(update, ctx, data[1])
		}
	}
}

func handleLvl3(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			delete(state.Data, "ActiveInput")
			HandleAddAds(update, ctx, "0")
		default:
			ActiveInput, exsist := state.Data["ActiveInput"].(ActiveInput)
			if exsist {
				HandleAddInput(update, ctx, strconv.Itoa(int(ActiveInput.ID)))
			}
		}

	}
	if update.Message != nil {
		state := context.GetUserState(userID, ctx)
		ActiveInput, exsist := state.Data["ActiveInput"].(ActiveInput)
		if exsist {
			HandleAddInput(update, ctx, strconv.Itoa(int(ActiveInput.ID)))
		}
	}
}
func handleLvl4(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	AdsPhoto, exsist := state.Data["AdsPhoto"].(AdsPhoto)
	if exsist {
		if AdsPhoto.MessageId != 0 {
			deleteMsg1 := tgbotapi.DeleteMessageConfig{
				ChatID:    userID,
				MessageID: AdsPhoto.MessageId,
			}
			AdsPhoto.MessageId = 0
			state.Data["AdsPhoto"] = AdsPhoto
			ctx.BotAPI.Send(deleteMsg1)
		}
		if update.CallbackQuery != nil {
			println(update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {
			case "back":
				AdsPhoto.ActivStep = 0
				AdsPhoto.IsEdit = true
				state.Data["AdsPhoto"] = AdsPhoto
				HandleAddAds(update, ctx, "0")
			default:
				HandleAddPhoto(update, ctx)
			}
		}

		if update.Message != nil {
			HandleAddPhoto(update, ctx)
		}
	}
}

func handelLvl5(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	AdsPhoto, exsist := state.Data["AdsPhoto"].(AdsPhoto)
	if exsist {
		if AdsPhoto.MessageId != 0 {
			deleteMsg1 := tgbotapi.DeleteMessageConfig{
				ChatID:    userID,
				MessageID: AdsPhoto.MessageId,
			}
			AdsPhoto.MessageId = 0
			state.Data["AdsPhoto"] = AdsPhoto
			ctx.BotAPI.Send(deleteMsg1)
		}
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "back":
				HandleAddAds(update, ctx, "0")
			}
		}
	}
}

func handelLvl6(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			HandleAddAds(update, ctx, "0")
		}
	}
}
