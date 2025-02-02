package profile

import (
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type transferType struct {
	ActiveStep uint
	Data       transferData
}

type transferData struct {
	RecipientID int64
	Amount      uint
}

func HandleDoPayment(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	var value string
	if update.Message != nil {
		userID = update.Message.Chat.ID
		value = update.Message.Text
		deleteMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg)
	} else {
		userID = update.CallbackQuery.From.ID
	}
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 6)

	var rows [][]tgbotapi.InlineKeyboardButton
	transfer, exist := state.Data["transfer"].(transferType)
	if !exist {
		transfer.ActiveStep = 0
		state.Data["transfer"] = transfer
	}

	switch transfer.ActiveStep {
	case 0:
		if update.CallbackQuery != nil && strings.Split(update.CallbackQuery.Data, "_")[0] == "Transfer" {
			transfer.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
			))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"1️⃣ Укажите ID пользователя, которому вы хотите перевести средства.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 1:
		if update.Message != nil {
			recipientID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
				))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Пожалуйста, укажите корректный ID пользователя (только цифры).",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			var count int64
			db.DB.Model(models.User{}).Where(models.User{TelegramID: int64(recipientID)}).Count(&count)
			if count == 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
				))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Пользователь не найден!",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			if userID == recipientID {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
				))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Вы не можете перевести самому себе!",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			transfer.Data.RecipientID = recipientID
			transfer.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
			))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"2️⃣ Укажите сумму, которую вы хотите перевести.\n\n💰 <i>Пример:</i> 500\n\n✍️ Пожалуйста, введите сумму перевода.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 2:
		if update.Message != nil {
			amount, err := strconv.ParseUint(value, 10, 64)
			if err != nil || amount <= 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
				))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Пожалуйста, укажите корректную сумму (например, 500).",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			transfer.Data.Amount = uint(amount)
			transfer.ActiveStep++
			var RecipientUser models.User
			db.DB.Where(models.User{TelegramID: int64(transfer.Data.RecipientID)}).First(&RecipientUser)
			var user models.User
			db.DB.Where("telegram_id = ?", userID).First(&user)
			if user.Balance < uint(amount) {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
				))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ На балансе недостаочно средств.",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", "confirm"),
			))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back"),
			))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"❗ Перевод на сумму: "+value+"\n\nПользователю: @"+RecipientUser.Username+" ("+RecipientUser.FirstName+" "+RecipientUser.LastName+")",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 3:
		if update.CallbackQuery != nil && update.CallbackQuery.Data == "confirm" {
			var user models.User
			db.DB.Where("telegram_id = ?", userID).First(&user)
			db.DB.Model(&models.User{}).
				Where("telegram_id = ?", userID).
				Updates(map[string]interface{}{
					"balance": user.Balance - uint(transfer.Data.Amount),
				})
			var RecipientUser models.User
			db.DB.Where(models.User{TelegramID: int64(transfer.Data.RecipientID)}).First(&RecipientUser)
			msg := tgbotapi.NewMessage(RecipientUser.TelegramID, "Перевод от пользователя"+user.Username+"\n\nНа сумму: "+strconv.Itoa(int(transfer.Data.Amount)))
			ctx.BotAPI.Send(msg)
			db.DB.Model(&models.User{}).
				Where(models.User{TelegramID: int64(transfer.Data.RecipientID)}).
				Updates(map[string]interface{}{
					"balance": RecipientUser.Balance + transfer.Data.Amount,
				})
			HandleProfile(update, ctx)
			delete(state.Data, "transfer")
		}
	}
	state.Data["transfer"] = transfer
}
