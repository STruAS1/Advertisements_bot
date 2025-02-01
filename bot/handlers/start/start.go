package start

import (
	"log"
	"strconv"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"tgbotBARAHOLKA/utilits"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStartCommand(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	if update.Message != nil {
		userID = update.Message.Chat.ID
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg1)

	} else {
		userID = update.CallbackQuery.From.ID
	}
	context.UpdateUserLevel(userID, ctx, 0)
	state := context.GetUserState(userID, ctx)
	var user models.User
	result := db.DB.Preload("Bans").Where("telegram_id = ?", userID).First(&user)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[0].ButtonText, "adsMenu"), tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[2].ButtonText, "profile"),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[1].ButtonText, "Docs"),
		},
	)
	if result.Error == nil {
		isBan, _ := user.IsBanned()
		if isBan {
			msg := tgbotapi.NewMessage(userID, utilits.FormatBanMessage(user))
			ctx.SendMessage(msg)
			return
		}
		if state.MessageID != 0 {
			msg := tgbotapi.NewEditMessageTextAndMarkup(userID, state.MessageID, config.GlobalSettings.Texts.MainText+"ㅤ", inlineKeyboard)
			ctx.BotAPI.Send(msg)

		} else {
			msg := tgbotapi.NewMessage(userID, config.GlobalSettings.Texts.MainText+"ㅤ")
			msg.ReplyMarkup = inlineKeyboard
			ctx.SendMessage(msg)
		}
		return
	}

	msg := tgbotapi.NewMessage(userID, "Привет! Для продолжения регистрации, отправьте мне свой номер телефона.")
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.KeyboardButton{
					Text:           "Отправить номер",
					RequestContact: true,
				},
			},
		},
		ResizeKeyboard: true,
	}

	msg.ReplyMarkup = keyboard
	context.UpdateUserLevel(userID, ctx, 1)
	ctx.SendMessage(msg)
}
func HandleVerificationRequest(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 0)
	CallBack := "Verification_" + strconv.Itoa(int(state.MessageID))
	verSufix := " (" + strconv.Itoa(int(config.GlobalSettings.VerificationCost)) + "₩)"
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Верификация"+verSufix, CallBack),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "StartMenu"),
		},
	)
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		"Пройти верификацию",
		inlineKeyboard,
	)
	ctx.BotAPI.Send(msg)
}

func HandlePhoneNumberRequest(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.Message.Chat.ID
	state := context.GetUserState(userID, ctx)
	deleteMsg1 := tgbotapi.DeleteMessageConfig{
		ChatID:    userID,
		MessageID: update.Message.MessageID,
	}
	ctx.BotAPI.Send(deleteMsg1)
	if update.Message.Contact != nil {
		userPhone := update.Message.Contact.PhoneNumber

		state.Data["phone_number"] = userPhone

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Выбрать город", "ChooseCity"),
			},
		)
		msg := tgbotapi.NewMessage(userID, "Для завершения регистрации, выберите город из списка")
		msg.ReplyMarkup = inlineKeyboard
		deleteMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: state.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg)
		ctx.SendMessage(msg)
		context.UpdateUserLevel(userID, ctx, 2)
	} else {
		deleteMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: state.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg)
		msg := tgbotapi.NewMessage(userID, "Пожалуйста, отправьте свой номер телефона.")
		keyboard := tgbotapi.ReplyKeyboardMarkup{
			Keyboard: [][]tgbotapi.KeyboardButton{
				{
					tgbotapi.KeyboardButton{
						Text:           "Отправить номер",
						RequestContact: true,
					},
				},
			},
			ResizeKeyboard: true,
		}

		msg.ReplyMarkup = keyboard
		ctx.SendMessage(msg)
	}
}

func HandleSubscriptionCheck(update *tgbotapi.Update, ctx *context.Context) {
	channelUsername := ctx.Config.Bot.ChannelId
	userID := update.CallbackQuery.From.ID

	config := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: channelUsername,
			UserID:             userID,
		},
	}

	member, err := ctx.BotAPI.GetChatMember(config)
	if err != nil {
		log.Printf("Ошибка при получении информации о пользователе: %v", err)
		return
	}
	state := context.GetUserState(userID, ctx)
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		userPhone := state.Data["phone_number"].(string)
		userCity := state.Data["CityTitle"].(string)
		newUser := models.User{
			TelegramID: update.CallbackQuery.From.ID,
			FirstName:  update.CallbackQuery.From.FirstName,
			Username:   update.CallbackQuery.From.UserName,
			LastName:   update.CallbackQuery.From.LastName,
			Phone:      userPhone,
			City:       userCity,
		}

		if err := db.DB.Create(&newUser).Error; err != nil {
			log.Printf("Ошибка при сохранении пользователя: %v", err)
			return
		}
		context.ClearAllUserData(userID, ctx)
		delete(state.Data, "phone_number")
		delete(state.Data, "CityTitle")
		context.UpdateUserLevel(userID, ctx, 0)
		HandleVerificationRequest(update, ctx)
	} else {
		alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "❌")
		alert.ShowAlert = false
		ctx.BotAPI.Request(alert)
	}
}
