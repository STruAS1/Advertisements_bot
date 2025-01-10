package ads

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"tgbotBARAHOLKA/utilits"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ActiveInput struct {
	ID           uint
	ActiveStep   uint
	Edit         bool
	Value        interface{}
	CitiesPages  map[uint]CitiesPage
	CurentPage   uint
	ActiveCities []CitiesRow
}

type CitiesRow struct {
	Id       uint
	Title    string
	IsActive bool
}

type CitiesPage struct {
	Cities []CitiesRow
}

func HandleAddInput(update *tgbotapi.Update, ctx *context.Context, InputID string) {
	var userID int64
	var value string
	var entities []tgbotapi.MessageEntity
	if update.Message != nil {
		userID = update.Message.Chat.ID
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		value = update.Message.Text
		entities = update.Message.Entities
		ctx.BotAPI.Send(deleteMsg1)

	} else {
		userID = update.CallbackQuery.From.ID
	}
	state := context.GetUserState(userID, ctx)

	context.UpdateUserLevel(userID, ctx, 3)
	inputIDInt, _ := strconv.Atoi(InputID)
	inputIDUint := uint(inputIDInt)
	adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
	var rows [][]tgbotapi.InlineKeyboardButton
	Input := adsInputs[inputIDUint]
	var Inputs models.AdvertisementInputs
	if state.Data["ActiveInput"] == nil {
		db.DB.Where(&models.AdvertisementInputs{ID: inputIDUint}).First(&Inputs)
		state.Data["ActiveInput"] = ActiveInput{
			ID:          inputIDUint,
			ActiveStep:  0,
			Value:       make(map[uint]string),
			Edit:        true,
			CitiesPages: make(map[uint]CitiesPage),
			CurentPage:  0,
		}
	} else {
		db.DB.Where(&models.AdvertisementInputs{ID: inputIDUint}).First(&Inputs)
	}
	ActiveInput := state.Data["ActiveInput"].(ActiveInput)
	if Input.Activate && ActiveInput.Edit {
		if update.CallbackQuery != nil {
			CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
			if CallbackQuery[0] == "AddInput" {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать", "Edit")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", "Delete")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "back")))
				text := "<b>" + Input.Name + "</b>: " + Input.Value
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.CallbackQuery.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
			if update.CallbackQuery.Data == "Delete" {
				switch ActiveInput.ActiveStep {
				case 0:
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", "Delete")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "back")))
					text := "❗ Вы уверены, что хотите удалить?"
					println(text)
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				case 1:
					adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
					Input.Activate = false
					Input.Value = ""
					adsInputs[inputIDUint] = Input
					HandleAddAds(update, ctx, "0")
					state := context.GetUserState(userID, ctx)
					delete(state.Data, "ActiveInput")
					return
				}
			}
			if update.CallbackQuery.Data == "Edit" {
				ActiveInput.Edit = false
				ActiveInput.ActiveStep = 0
				state.Data["ActiveInput"] = ActiveInput
			}
		}
	}
	switch Inputs.InputID {
	case 0:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Напишите любой короткий текст до 150 символов</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.Message != nil {
				if utf8.RuneCountInString(value) > 150 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Превышен лимит символов</b>` +
						"\n\n<i>❔Напишите любой короткий текст до 150 символов</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else {
					formatetText := utilits.ApplyFormatting(value, entities)
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = formatetText
					ActiveInput.Value = valueMap
					ActiveInput.ActiveStep = 2
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
					text := ("<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + formatetText)
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 2:
			if update.CallbackQuery.Data == "Save" {
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[0]
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 1:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Напишите любой текст до 2000 символов</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.Message != nil {
				if utf8.RuneCountInString(value) > 2000 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Превышен лимит символов</b>` +
						"\n\n<i>❔Напишите любой текст до 2000 символов</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					formatetText := utilits.ApplyFormatting(value, entities)
					valueMap[0] = formatetText
					ActiveInput.Value = valueMap
					ActiveInput.ActiveStep = 2
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
					text := ("<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>:\n" + formatetText)
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 2:
			if update.CallbackQuery.Data == "Save" {
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[0]
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 2:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("💳 Разовая", "OneTime")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔁 Регулярная", "Recurring")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип оплаты</i>" +
						"\n\n<blockquote><i>💳 Разовая: 100₩\n🔁 Регулярная: 10₩/Час</i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Recurring" {
					ActiveInput.ActiveStep = 2
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Час", "time_Час"), tgbotapi.NewInlineKeyboardButtonData("День", "time_День")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Неделя", "time_Неделя")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Месяц", "time_Месяц"), tgbotapi.NewInlineKeyboardButtonData("Год", "time_Год")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите план оплаты</i>" +
						"\n\n<blockquote><i>10₩/Час\n100₩/День\n700₩/Неделя\n3 000₩/Месяц\n36 500₩/Год</i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/<code>x</code>")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else if update.CallbackQuery.Data == "OneTime" {
					ActiveInput.ActiveStep = 3
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = "₩"
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📌 Фикс.Цена", "Fix"), tgbotapi.NewInlineKeyboardButtonData("🔄 Приблизительная цена", "Approximate")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📉📈 Диапазон", "Range")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип ценового ввода</i>" +
						"\n\n<blockquote><i><b>📌 Фикс.Цена</b>: 100₩\n\n<b>🔄 Приблизительная цена</b>: ~100₩\n\n<b>📉📈 Диапазон: 90-100₩</b></i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/<code>x</code>")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 2:
			if update.CallbackQuery != nil {
				Callback := strings.Split(update.CallbackQuery.Data, "_")
				if len(Callback) == 2 && Callback[0] == "time" {
					ActiveInput.ActiveStep = 3
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = "₩/" + Callback[1]
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📌 Фикс.Цена", "Fix"), tgbotapi.NewInlineKeyboardButtonData("🔄 Приблизительная цена", "Approximate")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📉📈 Диапазон", "Range")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип ценового ввода</i>" +
						"\n\n<blockquote><i><b>📌 Фикс.Цена</b>: 100₩/" + Callback[1] + "\n\n<b>🔄 Приблизительная цена</b>: ~100₩/" + Callback[1] + "\n\n<b>📉📈 Диапазон: 90-100₩/" + Callback[1] + "</b></i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/" + Callback[1])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 3:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Fix" || update.CallbackQuery.Data == "Approximate" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[1] = ""
					if update.CallbackQuery.Data == "Approximate" {
						valueMap[1] += "~"
					}
					ActiveInput.ActiveStep = 5
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите цену</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x</code>" + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else if update.CallbackQuery.Data == "Range" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[1] = ""
					ActiveInput.ActiveStep = 4
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите минимальную стоимость</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x - x</code>" + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 4:
			if update.Message != nil {
				price := strings.ReplaceAll(value, " ", "")
				priceFloat, err := strconv.ParseFloat(price, 64)
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
				if priceFloat > 10000000000000 {
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

				ActiveInput.ActiveStep = 5
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[1] = utilits.FormatFloatWithSpaces(priceFloat) + " - "
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<i>❔Введите максимальную стоимость</i>" +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x</code>" + valueMap[0])
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 5:
			if update.Message != nil {
				price := strings.ReplaceAll(value, " ", "")
				priceFloat, err := strconv.ParseFloat(price, 64)
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
				if priceFloat > 10000000000000 {
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

				ActiveInput.ActiveStep = 6
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[1] += utilits.FormatFloatWithSpaces(priceFloat)
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := ("\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + " " + valueMap[0])
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 6:
			if update.CallbackQuery.Data == "Save" {
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[1] + " " + valueMap[0]
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 3:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⚡ Сдельная", "OneTime")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔁 Регулярная", "Recurring")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип оплаты</i>" +
						"\n\n<blockquote><i><b>⚡ Сдельная</b>: 10₩/Шт\n<b>🔁 Регулярная</b>: 10₩/Час</i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.CallbackQuery.Data == "OneTime" {
				ActiveInput.ActiveStep = 2
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<i>❔Напишите идиницу измерения</i>" +
					"\n\n<blockquote><i><b>✅Пример:</b> Шт, М², Кг</i></blockquote>" +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/<code>x</code>")
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.CallbackQuery.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			} else if update.CallbackQuery.Data == "Recurring" {
				ActiveInput.ActiveStep = 3
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Час", "time_Час"), tgbotapi.NewInlineKeyboardButtonData("День", "time_День")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Неделя", "time_Неделя")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Месяц", "time_Месяц"), tgbotapi.NewInlineKeyboardButtonData("Год", "time_Год")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<i>❔Выберите план оплаты</i>" +
					"\n\n<blockquote><i>10₩/Час\n100₩/День\n700₩/Неделя\n3 000₩/Месяц\n36 500₩/Год</i></blockquote>" +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/<code>x</code>")
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.CallbackQuery.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 2:
			if update.Message != nil {
				value = strings.ReplaceAll(value, " ", "")
				if utf8.RuneCountInString(value) > 10 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := "<b>❗️ Превышен лимит символов!</b>"
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] += "₩/" + value
					ActiveInput.ActiveStep = 4
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📌 Фикс.Цена", "Fix"), tgbotapi.NewInlineKeyboardButtonData("🔄 Приблизительная цена", "Approximate")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📉📈 Диапазон", "Range")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип ценового ввода</i>" +
						"\n\n<blockquote><i><b>📌 Фикс.Цена</b>: 100" + valueMap[0] + "\n\n<b>🔄 Приблизительная цена</b>: ~100" + valueMap[0] + "\n\n<b>📉📈 Диапазон: 90-100" + valueMap[0] + "</b></i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>x</code> " + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 3:
			if update.CallbackQuery != nil {
				Callback := strings.Split(update.CallbackQuery.Data, "_")
				if len(Callback) == 2 && Callback[0] == "time" {
					ActiveInput.ActiveStep = 4
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = "₩/" + Callback[1]
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📌 Фикс.Цена", "Fix"), tgbotapi.NewInlineKeyboardButtonData("🔄 Приблизительная цена", "Approximate")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📉📈 Диапазон", "Range")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип ценового ввода</i>" +
						"\n\n<blockquote><i><b>📌 Фикс.Цена</b>: 100₩/" + Callback[1] + "\n\n<b>🔄 Приблизительная цена</b>: ~100₩/" + Callback[1] + "\n\n<b>📉📈 Диапазон: 90-100₩/" + Callback[1] + "</b></i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>x</code> ₩/" + Callback[1])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 4:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Fix" || update.CallbackQuery.Data == "Approximate" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[1] = ""
					if update.CallbackQuery.Data == "Approximate" {
						valueMap[1] += "~"
					}
					ActiveInput.ActiveStep = 6
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите зарплату</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x</code> " + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else if update.CallbackQuery.Data == "Range" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[1] = ""
					ActiveInput.ActiveStep = 5
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите минимальную зарплату</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x - x</code> " + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 5:
			if update.Message != nil {
				price := strings.ReplaceAll(value, " ", "")
				priceFloat, err := strconv.ParseFloat(price, 64)
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
				if priceFloat > 10000000000000 {
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

				ActiveInput.ActiveStep = 6
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[1] = utilits.FormatFloatWithSpaces(priceFloat) + " - "
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<i>❔Введите максимальную зарпалту</i>" +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + "<code>x</code>" + valueMap[0])
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 6:
			if update.Message != nil {
				price := strings.ReplaceAll(value, " ", "")
				priceFloat, err := strconv.ParseFloat(price, 64)
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
				if priceFloat > 10000000000000 {
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

				ActiveInput.ActiveStep = 7
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[1] += utilits.FormatFloatWithSpaces(priceFloat)
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := ("\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[1] + " " + valueMap[0])
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 7:
			if update.CallbackQuery.Data == "Save" {
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[1] + " " + valueMap[0]
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 4:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📌 Точное время", "Fix"), tgbotapi.NewInlineKeyboardButtonData("🔄 Приблизительное время", "Approximate")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📉📈 Временной диапазон", "Range")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите тип графика</i>" +
						"\n\n<blockquote><i><b>📌 Точное время</b>: 18:00\n\n<b>🔄 Приблизительное время</b>: ~18:00\n\n<b>📉📈 Временной диапазон: 18:00-19:00</b></i></blockquote>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Fix" || update.CallbackQuery.Data == "Approximate" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = ""
					if update.CallbackQuery.Data == "Approximate" {
						valueMap[0] += "~"
					}
					ActiveInput.ActiveStep = 3
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите время в формате <b>HH:mm</b></i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[0] + "<code>HH:mm</code>")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else if update.CallbackQuery.Data == "Range" {
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = ""
					ActiveInput.ActiveStep = 2
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Введите время начала в формате <b>HH:mm</b></i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: <code>HH:mm - HH:mm</code>")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 2:
			if update.Message != nil {
				time := strings.ReplaceAll(value, " ", "")
				regex := `^(?:[01]\d|2[0-3]):[0-5]\d$`
				matched, _ := regexp.MatchString(regex, time)
				if !matched {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := "❗️Ведите время в правильном формате <b>HH:mm</b>!"
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					ctx.BotAPI.Send(msg)
					return
				}

				ActiveInput.ActiveStep = 3
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[0] += time + " - "
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<i>❔Введите время окончания в формате <b>HH:mm</b></i>" +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[0] + "<code>HH:mm</code>")
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 3:
			if update.Message != nil {
				time := strings.ReplaceAll(value, " ", "")
				regex := `^(?:[01]\d|2[0-3]):[0-5]\d$`
				matched, _ := regexp.MatchString(regex, time)
				if !matched {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := "❗️Ведите время в правильном формате <b>HH:mm</b>!"
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					ctx.BotAPI.Send(msg)
					return
				}

				ActiveInput.ActiveStep = 4
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				valueMap[0] += time
				state.Data["ActiveInput"] = ActiveInput
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
					"\n\n<b>🔎Предварительный просмотр</b>" +
					"\n\n<b>" + Input.Name + "</b>: " + valueMap[0])
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.Message.Chat.ID,
					state.MessageID,
					text,
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				msg.ParseMode = "HTML"
				ctx.BotAPI.Send(msg)
				return
			}
		case 4:
			if update.CallbackQuery.Data == "Save" {
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[0]
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 5:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" {
					ActiveInput.ActiveStep = 1
					state.Data["ActiveInput"] = ActiveInput
					options := strings.Split(Inputs.Options, " ")
					for i := 0; i < len(options); i += 2 {
						var row []tgbotapi.InlineKeyboardButton
						button1 := tgbotapi.NewInlineKeyboardButtonData(options[i], "OPTIONS_"+options[i])
						row = append(row, button1)

						if i+1 < len(options) {
							button2 := tgbotapi.NewInlineKeyboardButtonData(options[i+1], "OPTIONS_"+options[i+1])
							row = append(row, button2)
						}

						rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
					}
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<i>❔Выберите любое значение из списка ниже</i>" +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: ")
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				fmt.Println(update.CallbackQuery.Data)
				if CallbackQuery[0] == "OPTIONS" {
					ActiveInput.ActiveStep = 2
					valueMap, _ := ActiveInput.Value.(map[uint]string)
					valueMap[0] = CallbackQuery[1]
					state.Data["ActiveInput"] = ActiveInput
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					text := (`<b>❗️Поле: "` + Input.Name + `"</b>` +
						"\n\n<b>🔎Предварительный просмотр</b>" +
						"\n\n<b>" + Input.Name + "</b>: " + valueMap[0])
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.From.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 2:
			if update.CallbackQuery.Data == "Save" {
				Input.Activate = true
				valueMap, _ := ActiveInput.Value.(map[uint]string)
				Input.Value = valueMap[0]
				adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
				adsInputs[inputIDUint] = Input
				HandleAddAds(update, ctx, "0")
				state := context.GetUserState(userID, ctx)
				delete(state.Data, "ActiveInput")
				return
			}
		}
	case 6:
		switch ActiveInput.ActiveStep {
		case 0:
			if update.CallbackQuery != nil {
				CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
				if CallbackQuery[0] == "AddInput" || update.CallbackQuery.Data == "Edit" || update.CallbackQuery.Data == "nextCity" || update.CallbackQuery.Data == "backCity" || update.CallbackQuery.Data == "search" || update.CallbackQuery.Data == "BackToList" || CallbackQuery[0] == "City" || update.CallbackQuery.Data == "Save" || update.CallbackQuery.Data == "menuCityInfo" {
					if update.CallbackQuery.Data == "nextCity" && len(ActiveInput.CitiesPages)-1 != int(ActiveInput.CurentPage) {
						ActiveInput.CurentPage++
						state.Data["ActiveInput"] = ActiveInput
					}
					if update.CallbackQuery.Data == "backCity" && int(ActiveInput.CurentPage) != 0 {
						ActiveInput.CurentPage--
						state.Data["ActiveInput"] = ActiveInput
					}
					if update.CallbackQuery.Data == "menuCityInfo" {
						alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Выберете город из списка")
						alert.ShowAlert = false
						ctx.BotAPI.Request(alert)
						return
					}
					if update.CallbackQuery.Data == "Save" {
						var textActiveCities string = ""
						if len(ActiveInput.ActiveCities) != 0 {
							for i := range ActiveInput.ActiveCities {
								textActiveCities += ActiveInput.ActiveCities[i].Title
								if i+1 < len(ActiveInput.ActiveCities) {
									textActiveCities += ", "
								} else {
									textActiveCities += "."
								}
							}
							Input.Activate = true
							Input.Value = textActiveCities
							adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
							adsInputs[inputIDUint] = Input
							HandleAddAds(update, ctx, "0")
							state := context.GetUserState(userID, ctx)
							delete(state.Data, "ActiveInput")
							return
						} else {
							alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Вы не можете сохранить, пока не выбрали ни одного горада!")
							ctx.BotAPI.Request(alert)
						}

					}
					ActiveInput.ActiveStep = 0
					if update.CallbackQuery.Data == "search" {
						var textActiveCities string = ""
						if len(ActiveInput.ActiveCities) != 0 {
							textActiveCities += "<b>" + Input.Name + "</b>: "
							for i := range ActiveInput.ActiveCities {
								textActiveCities += ActiveInput.ActiveCities[i].Title
								if i+1 < len(ActiveInput.ActiveCities) {
									textActiveCities += ", "
								} else {
									textActiveCities += "."
								}
							}
						}
						ActiveInput.ActiveStep = 1
						state.Data["ActiveInput"] = ActiveInput
						text := "Введите название города" + "\n" + textActiveCities
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад к списку", "BackToList")))
						msg := tgbotapi.NewEditMessageTextAndMarkup(
							update.CallbackQuery.Message.Chat.ID,
							state.MessageID,
							text,
							tgbotapi.NewInlineKeyboardMarkup(rows...),
						)
						msg.ParseMode = "HTML"
						ctx.BotAPI.Send(msg)
						return
					}
					if len(ActiveInput.CitiesPages) == 0 {
						pageSize := 10
						var cities []models.Cities
						db.DB.Order("title ASC").Find(&cities)
						for page := 0; page < (len(cities)+pageSize-1)/pageSize; page++ {
							var _Cities []CitiesRow
							start := page * pageSize
							end := start + pageSize
							if end > len(cities) {
								end = len(cities)
							}

							for _, city := range cities[start:end] {
								_Cities = append(_Cities, CitiesRow{Id: city.ID, Title: city.Title, IsActive: false})
							}

							ActiveInput.CitiesPages[uint(page)] = CitiesPage{Cities: _Cities}

						}

						state.Data["ActiveInput"] = ActiveInput
					}
					if CallbackQuery[0] == "City" {
						cytyArrayID, _ := strconv.Atoi(CallbackQuery[3])
						pageID, _ := strconv.Atoi(CallbackQuery[2])
						ActiveInput.CurentPage = uint(pageID)
						if ActiveInput.CitiesPages[uint(pageID)].Cities[cytyArrayID].IsActive {
							ActiveInput.CitiesPages[uint(pageID)].Cities[cytyArrayID].IsActive = false
							cityToRemove := ActiveInput.CitiesPages[uint(pageID)].Cities[cytyArrayID]
							for i, city := range ActiveInput.ActiveCities {
								if city.Id == cityToRemove.Id {
									ActiveInput.ActiveCities = append(ActiveInput.ActiveCities[:i], ActiveInput.ActiveCities[i+1:]...)
									break
								}
							}
						} else {
							if len(ActiveInput.ActiveCities) >= 10 {
								alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Вы не можете добавить больше 10 городов!")
								ctx.BotAPI.Request(alert)
							} else {
								ActiveCity := ActiveInput.CitiesPages[uint(pageID)].Cities[cytyArrayID]
								ActiveInput.ActiveCities = append(ActiveInput.ActiveCities, ActiveCity)
								ActiveInput.CitiesPages[uint(pageID)].Cities[cytyArrayID].IsActive = true
							}
						}
						state.Data["ActiveInput"] = ActiveInput
					}
					currentPage := ActiveInput.CurentPage
					page := ActiveInput.CitiesPages[currentPage].Cities
					for i := 0; i < len(page); i += 2 {
						var titleI string = page[i].Title
						if page[i].IsActive {
							titleI += " ✅"
						}
						if i+1 < len(page) {
							var titleI1 string = page[i+1].Title
							if page[i+1].IsActive {
								titleI1 += " ✅"
							}
							rows = append(rows, tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData(titleI, "City_"+strconv.Itoa(int(page[i].Id))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i))),
								tgbotapi.NewInlineKeyboardButtonData(titleI1, "City_"+strconv.Itoa(int(page[i+1].Id))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i+1))),
							))

						} else {
							rows = append(rows, tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData(titleI, "City_"+strconv.Itoa(int(page[i].Id))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i))),
							))
						}
					}
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✩✩✩ ", "menuCityInfo")))
					if len(ActiveInput.CitiesPages)-1 > int(currentPage) && currentPage != 0 {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "backCity"), tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search"), tgbotapi.NewInlineKeyboardButtonData("Дальше »", "nextCity")))
					} else if len(ActiveInput.CitiesPages)-1 > int(currentPage) && currentPage == 0 {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search"), tgbotapi.NewInlineKeyboardButtonData("Дальше »", "nextCity")))
					} else if len(ActiveInput.CitiesPages)-1 == int(currentPage) {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "backCity"), tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search")))
					}
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					var textActiveCities string = ""
					if len(ActiveInput.ActiveCities) != 0 {
						textActiveCities += "<b>" + Input.Name + "</b>: "
						for i := range ActiveInput.ActiveCities {
							textActiveCities += ActiveInput.ActiveCities[i].Title
							if i+1 < len(ActiveInput.ActiveCities) {
								textActiveCities += ", "
							} else {
								textActiveCities += "."
							}
						}
					}
					text := "🏙️ Выберите город из списка!" + "\n" + textActiveCities
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.Message != nil {
				var matchedCities []struct {
					City      CitiesRow
					PageID    uint
					CityIndex int
				}
				var otherCities []struct {
					City      CitiesRow
					PageID    uint
					CityIndex int
				}

				for pageID, page := range ActiveInput.CitiesPages {
					for cityIndex, city := range page.Cities {
						if strings.Contains(strings.ToLower(city.Title), strings.ToLower(value)) {
							cityData := struct {
								City      CitiesRow
								PageID    uint
								CityIndex int
							}{
								City:      city,
								PageID:    pageID,
								CityIndex: cityIndex,
							}

							if strings.HasPrefix(strings.ToLower(city.Title), strings.ToLower(value)) {
								matchedCities = append(matchedCities, cityData)
							} else {
								otherCities = append(otherCities, cityData)
							}
						}
					}
				}

				sortedCities := append(matchedCities, otherCities...)

				if len(sortedCities) > 10 {
					sortedCities = sortedCities[:10]
				}
				for _, cityData := range sortedCities {
					var title string = cityData.City.Title
					if cityData.City.IsActive {
						title += " ✅"
					}
					data := fmt.Sprintf("City_%d_%d_%d", cityData.City.Id, cityData.PageID, cityData.CityIndex)

					if len(rows) == 0 || len(rows[len(rows)-1]) == 2 {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(title, data),
						))
					} else {
						rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData(title, data))
					}
				}
				var textActiveCities string = ""
				if len(ActiveInput.ActiveCities) != 0 {
					textActiveCities += "<b>" + Input.Name + "</b>: "
					for i := range ActiveInput.ActiveCities {
						textActiveCities += ActiveInput.ActiveCities[i].Title
						if i+1 < len(ActiveInput.ActiveCities) {
							textActiveCities += ", "
						} else {
							textActiveCities += "."
						}
					}
				}
				if len(sortedCities) == 0 {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Совпадений не найдено", "NoResult")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад к списку", "BackToList")))
					text := "Введите название Города" + "\n" + textActiveCities
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				} else {
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад к списку", "BackToList")))
					text := "Введите название Города" + "\n" + textActiveCities
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.Message.Chat.ID,
						state.MessageID,
						text,
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(msg)
					return
				}
			}
			if update.CallbackQuery != nil {
				callbackQuery := update.CallbackQuery
				data := callbackQuery.Data
				if data == "NoResult" {
					callback := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "Город не найден!")
					callback.ShowAlert = false
					ctx.BotAPI.Request(callback)
				}
			}
		}
	}
}
