package verification

import (
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type verificationType struct {
	ActiveStep uint
	Data       verificationData
}

type verificationData struct {
	FIO            FIO
	VisaType       string
	Services       string
	CardIdFileID   string
	DocumentFileID string
}

type FIO struct {
	FirstName  string
	LastName   string
	Patronymic string
}

func HandleBackToStartMenu(ctx *context.Context, userID int64) {
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 0)
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Главное меню", "StartMenu")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		"⏳ Спасибо за предоставленные данные! Ожидайте ручной проверки.",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}

func HandleVerification(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	var value string
	var photoID string
	var user models.User
	db.DB.Where("telegram_id = ?", userID).First(&user)
	if update.Message != nil {
		userID = update.Message.Chat.ID
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		if update.Message.Photo != nil {
			photoID = update.Message.Photo[len(update.Message.Photo)-1].FileID
		}
		value = update.Message.Text
		ctx.BotAPI.Send(deleteMsg1)

	} else {
		userID = update.CallbackQuery.From.ID
	}

	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 1)

	var rows [][]tgbotapi.InlineKeyboardButton
	verification, exist := state.Data["verification"].(verificationType)
	if !exist {
		verification.ActiveStep = 0
		state.Data["verification"] = verification
	}

	switch verification.ActiveStep {
	case 0:
		if update.CallbackQuery != nil && strings.Split(update.CallbackQuery.Data, "_")[0] == "Verification" {
			if user.Balance < config.GlobalSettings.VerificationCost {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"Недостаточно средств",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				ctx.BotAPI.Send(msg)
			}
			verification.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"1️⃣ Введите ФИО как в удостоверении личности (ID-card)\n\n📋 <i>Пример:</i> 홍길동 / Hong Gil Dong\n\n✍️ Пожалуйста, введите ваши ФИО, как они указаны в вашей ID-карте.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 1:
		if update.Message != nil {
			if value == "" {
				msg := tgbotapi.NewMessage(userID, "⚠️ Пожалуйста, введите ФИО как в удостоверении личности.")
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			// FIO := strings.Split(value, " ")
			verification.Data.FIO.LastName = value
			// verification.Data.FIO.FirstName = FIO[1]
			// verification.Data.FIO.Patronymic = FIO[2]
			verification.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"2️⃣ Укажите тип визы или статус пребывания\n\n🌍 <i>Пример:</i> Гражданство Кореи, F-5 (ВНЖ), F-4, F-6, F-2, F-1...\n\n🔑 Напишите ваш текущий тип визы или статус.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 2:
		if update.Message != nil {
			if value == "" {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Пожалуйста, укажите тип визы или статус пребывания.",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			verification.Data.VisaType = value
			verification.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"3️⃣ Перечислите услуги, которые вы предоставляете\n\n💼 <i>Пример:</i> Репетиторство, перевод документов, юридические консультации\n\n📜 Введите список через запятую.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 3:
		if update.Message != nil {
			if value == "" {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					userID,
					state.MessageID,
					"⚠️ Пожалуйста, укажите перечень услуг.",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			verification.Data.Services = value
			verification.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"4️⃣ Загрузите фото удостоверения личности (ID-card)\n\n📸 Сделайте фото или загрузите скан вашей ID-карты.\n\n🔒 <b>Ваши данные защищены!</b>",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 4:
		if update.Message.Photo != nil {
			verification.Data.CardIdFileID = photoID
			verification.ActiveStep++
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить", "skip")))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"5️⃣ Загрузите документ, подтверждающий компетенции (опционально)\n\n📚 <i>Пример:</i> Диплом, лицензия, свидетельство о регистрации компании\n\n🌟 Если у вас нет таких документов, вы можете пропустить этот шаг.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена", "back")))
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				userID,
				state.MessageID,
				"⚠️ Пожалуйста, загрузите фото удостоверения личности.",
				tgbotapi.NewInlineKeyboardMarkup(rows...),
			)
			msg.ParseMode = "HTML"
			ctx.BotAPI.Send(msg)
		}
	case 5:
		if update.Message != nil && update.Message.Photo != nil {
			verification.Data.DocumentFileID = photoID
		}
		delete(state.Data, "verification")
		HandleBackToStartMenu(ctx, userID)
		if err := db.DB.Model(&models.User{}).
			Where("id = ?", uint(user.ID)).
			Updates(map[string]interface{}{
				"verification": true,
				"balance":      user.Balance - config.GlobalSettings.VerificationCost,
			}).Error; err != nil {
			return
		}
		NewAplication := models.VerificationAplication{
			UserID:         user.ID,
			FirstName:      verification.Data.FIO.FirstName,
			LastName:       verification.Data.FIO.LastName,
			Patronymic:     verification.Data.FIO.Patronymic,
			VisaType:       verification.Data.VisaType,
			CardIdFileID:   verification.Data.CardIdFileID,
			DocumentFileID: verification.Data.DocumentFileID,
			Services:       verification.Data.Services,
		}
		db.DB.Create(&NewAplication)
		return

	}
	state.Data["verification"] = verification
}
