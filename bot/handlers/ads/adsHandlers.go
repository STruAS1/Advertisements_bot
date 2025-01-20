package ads

import (
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	switch state.Level {
	case 1:
		handleLvl1(update, ctx)
	case 2:
		handleLvl2(update, ctx)
	case 3:
		handleLvl3(update, ctx, userID)
	case 4:
		handleLvl4(update, ctx, userID)
	case 5:
		handelLvl5(update, ctx, userID)
	case 6:
		handelLvl6(update, ctx)
	case 7:
		handelLvl7(update, ctx, userID)
	case 8:
		handelLvl8(update, ctx, userID)
	case 9:
		handelLvl9(update, ctx, userID)
	case 10:
		handleLvl10(update, ctx, userID)
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
			HandleAddAds(update, ctx, data[1], false)
		}
	}
}

func handleLvl2(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		if len(data) == 1 {
			switch data[0] {
			case "back":
				HandleBackAds(update, ctx)
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
		data := strings.Split(update.CallbackQuery.Data, "_")
		switch data[0] {
		case "back":
			delete(state.Data, "ActiveInput")
			HandleAddAds(update, ctx, "0", false)
		case "BackToList":
			ActiveInput, exsist := state.Data["ActiveInput"].(ActiveInput)
			if exsist {
				ActiveInput.ActiveStep = 0
				state.Data["ActiveInput"] = ActiveInput
				HandleAddInput(update, ctx, strconv.Itoa(int(ActiveInput.ID)))
			}
		case "City":
			ActiveInput, exsist := state.Data["ActiveInput"].(ActiveInput)
			if exsist {
				ActiveInput.ActiveStep = 0
				state.Data["ActiveInput"] = ActiveInput
				HandleAddInput(update, ctx, strconv.Itoa(int(ActiveInput.ID)))
			}
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
				HandleAddAds(update, ctx, "0", false)
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
				HandleAddAds(update, ctx, "0", false)
			}
		}
	}
}

func handelLvl6(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			HandleAddAds(update, ctx, "0", false)
		}
	}
}

func handelLvl7(update *tgbotapi.Update, ctx *context.Context, userId int64) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			HandleAddAds(update, ctx, "0", false)
		case "Delete":
			state := context.GetUserState(userId, ctx)
			delete(state.Data, "AdsInputs")
			delete(state.Data, "AdsPhoto")
			delete(state.Data, "ActivType")
			delete(state.Data, "SkipTimerCoast")
			HandleMenu(update, ctx)
		}

	}
}

func handelLvl8(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		state := context.GetUserState(userID, ctx)
		switch strings.Split(update.CallbackQuery.Data, "_")[0] {
		case "back":
			delete(state.Data, "AdsHistoryPage")
			delete(state.Data, "AdsHistory")
			HandleMenu(update, ctx)
		case "backAds":
			pages := state.Data["AdsHistory"].(map[uint]pageHistoryAds)
			ActivePage := state.Data["AdsHistoryPage"].(uint)
			if len(pages)-1 != int(ActivePage) {
				ActivePage--
				state.Data["AdsHistoryPage"] = ActivePage
				HandleSelectADSHistory(update, ctx)
			}
		case "nextAds":
			ActivePage := state.Data["AdsHistoryPage"].(uint)
			if int(ActivePage) != 0 {
				ActivePage++
				state.Data["AdsHistoryPage"] = ActivePage
				HandleSelectADSHistory(update, ctx)
			}
		case "Ad":
			HandleViwerAdsHistory(update, ctx)
		}
	}
}

func handelLvl9(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			state := context.GetUserState(userID, ctx)
			if state.Data["MessageIdPhoto"] == 0 {
				delete(state.Data, "MessageIdPhoto")
				HandleSelectADSHistory(update, ctx)
			} else {
				messageId := state.Data["MessageIdPhoto"].(int)
				deleteMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    userID,
					MessageID: messageId,
				}
				ctx.BotAPI.Send(deleteMsg)
				delete(state.Data, "MessageIdPhoto")
				HandleSelectADSHistory(update, ctx)
			}
		}
	}
}

func handleLvl10(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		if len(data) == 1 {
			switch data[0] {
			case "back":
				HandleSelectADS(update, ctx)
			}
		} else if len(data) == 2 && data[0] == "buy" {
			var User models.User
			db.DB.Where(&models.User{TelegramID: userID}).First(&User)
			if User.Balance < config.GlobalSettings.Ads.CostLimit {
				message := "Недостаточно средств на балансе!"
				alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, message)
				alert.ShowAlert = false
				ctx.BotAPI.Request(alert)
				return
			}
			state := context.GetUserState(userID, ctx)
			state.Data["SkipTimerCoast"] = config.GlobalSettings.Ads.CostLimit
			HandleAddAds(update, ctx, data[1], true)
		}
	}
}
