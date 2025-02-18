package start

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
		if update.Message != nil {
			HandlePhoneNumberRequest(update, ctx)
		}
	case 2:
		handleLvl2(update, ctx)
	case 3:
		handelLvl3(update, ctx)
	case 4:
		handleLvl4(update, ctx, userID)
	case 5:
		handelLvl5(update, ctx, userID)
	}

}

func handleLvl2(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "ChooseCity":
			ChooseCityHandler(update, ctx)
		case "Chek_sub":
			HandleSubscriptionCheck(update, ctx)
		}

	}
}

func handleLvl4(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	if update.CallbackQuery != nil {
		switch strings.Split(update.CallbackQuery.Data, "_")[0] {
		case "City":
			CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
			cytyArrayID, _ := strconv.Atoi(CallbackQuery[3])
			pageID, _ := strconv.Atoi(CallbackQuery[2])
			ActiveChooseCity, _ := state.Data["ActiveChooseCity"].(ActiveChooseCityType)
			state.Data["CityTitle"] = ActiveChooseCity.CitiesPages[uint(pageID)].Cities[cytyArrayID].Title
			delete(state.Data, "ActiveChooseCity")
			var user models.User
			result := db.DB.Where("telegram_id = ?", userID).First(&user)
			if result.Error == nil {
				if err := db.DB.Model(&models.User{}).
					Where("id = ?", uint(user.ID)).
					Updates(map[string]interface{}{
						"city": state.Data["CityTitle"].(string),
					}).Error; err != nil {
					return
				}
				HandleStartCommand(update, ctx)
				return
			} else {
				channelUsername := ctx.Config.Bot.ChannelId
				channelUsername = strings.TrimPrefix(channelUsername, "@")
				url := "https://t.me/" + channelUsername

				inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					[]tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonURL(config.GlobalSettings.Buttons[3].ButtonText, url),
					},
					[]tgbotapi.InlineKeyboardButton{
						tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[4].ButtonText, "Chek_sub"),
					},
				)
				msg := tgbotapi.NewMessage(userID, "Для завершения регистрации, пожалуйста, подпишитесь на канал.")
				msg.ReplyMarkup = inlineKeyboard
				deleteMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    userID,
					MessageID: state.MessageID,
				}
				ctx.BotAPI.Send(deleteMsg)
				ctx.SendMessage(msg)
				context.UpdateUserLevel(userID, ctx, 2)
			}
		default:
			ChooseCityHandler(update, ctx)
		}
	}
	if update.Message != nil {
		ChooseCityHandler(update, ctx)
	}
}
func handelLvl3(update *tgbotapi.Update, ctx *context.Context) {
	if update.CallbackQuery != nil {
		data := strings.Split(update.CallbackQuery.Data, "_")
		switch data[0] {
		case "back":
			HandleStartCommand(update, ctx)
		case "doc":
			if len(data) == 2 {
				Index, _ := strconv.Atoi(data[1])
				HandleDocs(update, ctx, Index)
			}
		}
	}
}

func handelLvl5(update *tgbotapi.Update, ctx *context.Context, userID int64) {
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			state := context.GetUserState(userID, ctx)
			MessageId, exist := state.Data["LastVideoMassgeID"].(int)
			if exist {
				deleteMsg1 := tgbotapi.DeleteMessageConfig{
					ChatID:    userID,
					MessageID: MessageId,
				}
				ctx.BotAPI.Send(deleteMsg1)
				delete(state.Data, "LastVideoMassgeID")
			}
			HandleSelectDocs(update, ctx)
		}
	}
}
