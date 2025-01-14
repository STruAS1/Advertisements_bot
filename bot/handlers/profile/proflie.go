package profile

import (
	"fmt"
	"strconv"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleProfile(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 0)
	var rows [][]tgbotapi.InlineKeyboardButton
	var user models.User
	db.DB.Where("telegram_id = ?", userID).First(&user)
	text := "<b>" + user.FirstName + " " + user.LastName + "</b>"
	text += "\n<b>ID</b>: <code>" + strconv.Itoa(int(user.ID)) + "</code>"
	text += "\n\n<b>Баланс</b>: " + strconv.Itoa(int(user.Balance))
	var CountOfAds int64
	db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID}).Count(&CountOfAds)
	var AprovedCounOFAds int64
	db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID, Status: 1}).Count(&AprovedCounOFAds)
	text += "\n\n<b>Всего объявлений</b>: " + strconv.Itoa(int(CountOfAds))
	text += "\n<b>Опубликовано объявлений</b>: " + strconv.Itoa(int(AprovedCounOFAds))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пополнить баланс", "+balance")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "StartMenu")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		text,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	_, err := ctx.BotAPI.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}
