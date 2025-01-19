package start

import (
	"log"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

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
	result := db.DB.Where("telegram_id = ?", userID).First(&user)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Объявление", "adsMenu"), tgbotapi.NewInlineKeyboardButtonData("Профиль", "profile"),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Обучение", "Docs"),
		},
	)
	if result.Error == nil {
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
		state.Data["Text_msg"] = "Для завершения регистрации, пожалуйста, подпишитесь на канал."
		channelUsername := ctx.Config.Bot.ChannelId
		channelUsername = strings.TrimPrefix(channelUsername, "@")
		url := "https://t.me/" + channelUsername

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonURL("Подписаться на канал", url),
			},
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Подписался", "cehk_sub"),
			},
		)
		msg := tgbotapi.NewMessage(userID, "Спасибо за номер. Теперь подпишитесь на канал.")
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

		newUser := models.User{
			TelegramID: update.CallbackQuery.From.ID,
			FirstName:  update.CallbackQuery.From.FirstName,
			Username:   update.CallbackQuery.From.UserName,
			LastName:   update.CallbackQuery.From.LastName,
			Phone:      userPhone,
		}

		if err := db.DB.Create(&newUser).Error; err != nil {
			log.Printf("Ошибка при сохранении пользователя: %v", err)
			return
		}
		context.ClearAllUserData(userID, ctx)

		context.UpdateUserLevel(userID, ctx, 0)
		HandleStartCommand(update, ctx)
	} else {
		alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "❌")
		alert.ShowAlert = false
		ctx.BotAPI.Request(alert)
	}
}
