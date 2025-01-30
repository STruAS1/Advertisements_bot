package profile

import (
	"fmt"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleProfile(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	if update.Message != nil {
		userID = update.Message.Chat.ID
		deleteMsg := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg)
	} else {
		userID = update.CallbackQuery.From.ID
	}
	seoulLocation, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		fmt.Println("Ошибка загрузки временной зоны:", err)
		return
	}

	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 0)
	var rows [][]tgbotapi.InlineKeyboardButton
	var user models.User
	db.DB.Where("telegram_id = ?", userID).First(&user)
	kstTime := user.CreatedAt.In(seoulLocation)
	formattedDate := kstTime.Format("2006.01.02")
	text := fmt.Sprintf(
		"<b>%s %s</b>\n<b>ID</b>: <code>%d</code>\n<b>Телефон</b>: <code>%s</code>\n<b>Населенный пункт</b>: <code>%s</code>\n<b>Дата регистрации</b>: <code>%s</code>\n<b>Баланс</b>: %d",
		user.FirstName, user.LastName, userID, user.Phone, user.City, formattedDate, user.Balance,
	)
	var CountOfAds int64
	db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID}).Count(&CountOfAds)
	var AprovedCounOFAds int64
	db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID, Status: 1}).Count(&AprovedCounOFAds)
	text += "\n\n<b>Всего объявлений</b>: " + strconv.Itoa(int(CountOfAds))
	text += "\n<b>Опубликовано объявлений</b>: " + strconv.Itoa(int(AprovedCounOFAds))
	var VerStatus string = "✅"
	var CallBack string = "nil"
	var verSufix string
	if !user.Verification {
		var verification models.VerificationAplication
		result := db.DB.Where(&models.VerificationAplication{UserID: user.ID}).
			Order("id DESC").
			First(&verification)
		if result.Error != nil {
			VerStatus = ""
			verSufix = " (" + strconv.Itoa(int(config.GlobalSettings.VerificationCost)) + "₩)"
			CallBack = "Verification_" + strconv.Itoa(int(state.MessageID))
		} else if verification.Status == 0 {
			VerStatus = "⏳"
		} else if verification.Status == 2 {
			VerStatus = "❌"
			verSufix = " (" + strconv.Itoa(int(config.GlobalSettings.VerificationCost)) + "₩)"
			CallBack = "Verification_" + strconv.Itoa(int(state.MessageID))
		}
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(VerStatus+"Верификация"+verSufix, CallBack)))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[6].ButtonText, "+balance")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[7].ButtonText, "Transfer")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[5].ButtonText, "StartMenu")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		text,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	_, err = ctx.BotAPI.Send(msg)
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
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[5].ButtonText, "back")))
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
	Metod   config.PaymentsMetod
	Amount  float64
	PhotoId string
}

func HandlePaymentEntryAmount(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 2)
	metodIndex, _ := strconv.Atoi(strings.Split(update.CallbackQuery.Data, "_")[1])
	state.Data["Payment"] = Payment{Metod: config.GlobalSettings.Payments.Metods[metodIndex], Amount: 0}
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[5].ButtonText, "back")))
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
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[14].ButtonText, "back")))
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
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[14].ButtonText, "back")))
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
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[14].ButtonText, "back")))
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
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[14].ButtonText, "back")))
	var text string = "<b>" + payment.Metod.Title + "</b>\n\n"
	text += payment.Metod.Discription + "\n\n<b><i>Реквизиты: </i></b>"
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
func HandleGetPhoto(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.Message.From.ID
	deleteMsg1 := tgbotapi.DeleteMessageConfig{
		ChatID:    userID,
		MessageID: update.Message.MessageID,
	}
	state := context.GetUserState(userID, ctx)
	ctx.BotAPI.Send(deleteMsg1)
	if update.Message != nil && update.Message.Photo != nil {
		payment := state.Data["Payment"].(Payment)
		photoID := update.Message.Photo[len(update.Message.Photo)-1].FileID
		payment.PhotoId = photoID
		state.Data["Payment"] = payment
		HandeleConfirmPayment(update, ctx)
		return
	} else {
		var rows [][]tgbotapi.InlineKeyboardButton
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[14].ButtonText, "back")))
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			userID,
			state.MessageID,
			"Отпарьте фотографию",
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		ctx.BotAPI.Send(msg)
	}
}

func HandeleConfirmPayment(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.Message.From.ID
	state := context.GetUserState(userID, ctx)
	payment := state.Data["Payment"].(Payment)
	var user models.User
	db.DB.Where("telegram_id = ?", userID).First(&user)
	newPayment := models.Payments{
		Metod:    "(" + payment.Metod.Title + ")" + payment.Metod.Cardnumber,
		Amount:   uint(payment.Amount),
		UserID:   uint(user.ID),
		PhotoUrl: payment.PhotoId,
		Status:   0,
	}
	db.DB.Create(&newPayment)
	delete(state.Data, "Payment")
	HandleProfile(update, ctx)
}
