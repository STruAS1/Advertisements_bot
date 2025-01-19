package profile

import (
	"fmt"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
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

func HandleSelectPaymentMetod(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 1)
	metods := config.GlobalSettings.Payments.Metods
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, metod := range metods {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(metod.Title, "payment_"+strconv.Itoa(int(i)))))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		"Выберите способ пополнения",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	_, err := ctx.BotAPI.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

type Payment struct {
	Metod  config.PaymentsMetod
	Amount float64
}

func HandlePaymentEntryAmount(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 2)
	metodIndex, _ := strconv.Atoi(strings.Split(update.CallbackQuery.Data, "_")[1])
	state.Data["Payment"] = Payment{Metod: config.GlobalSettings.Payments.Metods[metodIndex], Amount: 0}
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		"Введите сумму",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	_, err := ctx.BotAPI.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleShowMetods(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.Message.From.ID
	deleteMsg1 := tgbotapi.DeleteMessageConfig{
		ChatID:    userID,
		MessageID: update.Message.MessageID,
	}
	ctx.BotAPI.Send(deleteMsg1)
	state := context.GetUserState(userID, ctx)
	payment := state.Data["Payment"].(Payment)
	price := strings.ReplaceAll(update.Message.Text, " ", "")
	priceFloat, err := strconv.ParseFloat(price, 64)
	var rows [][]tgbotapi.InlineKeyboardButton
	if err != nil {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
		text := "❗️Введите коректное числовое значение!"
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			update.Message.Chat.ID,
			state.MessageID,
			text,
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		ctx.BotAPI.Send(msg)
		return
	}
	if priceFloat > float64(config.GlobalSettings.Payments.MaxAmount) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
		text := "❗️Вы ввели больше максимального значения!"
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			update.Message.Chat.ID,
			state.MessageID,
			text,
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		ctx.BotAPI.Send(msg)
		return
	}
	if priceFloat < float64(config.GlobalSettings.Payments.MinimalAmount) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
		text := "❗️Вы ввели меньше минимального значения!"
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			update.Message.Chat.ID,
			state.MessageID,
			text,
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		ctx.BotAPI.Send(msg)
		return
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить ", "confirm")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
	var text string = "<b>" + payment.Metod.Title + "</b>\n\n"
	text += payment.Metod.Title + "\n\n<b><i>Реквизиты:</i></b>"
	text += payment.Metod.Cardnumber
	payment.Amount = priceFloat
	state.Data["Payment"] = payment
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		text,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
	context.UpdateUserLevel(userID, ctx, 3)
}

func HandeleConfirmPayment(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	payment := state.Data["Payment"].(Payment)
	var user models.User
	db.DB.Where("telegram_id = ?", userID).First(&user)
	newPayment := models.Payments{
		Metod:  "(" + payment.Metod.Title + ")" + payment.Metod.Cardnumber,
		Amount: uint(payment.Amount),
		UserID: uint(user.ID),
		Status: 0,
	}
	db.DB.Create(&newPayment)
	delete(state.Data, "Payment")
	HandleProfile(update, ctx)
}
