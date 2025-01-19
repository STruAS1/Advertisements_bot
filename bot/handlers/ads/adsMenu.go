package ads

import (
	"sort"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdsInputs struct {
	ID       uint
	Optional bool
	Name     string
	Activate bool
	Value    string
	Priority uint
}

type AdsPhoto struct {
	ID        string
	IDpre     string
	Activate  bool
	MessageId int
	ActivStep int
	IsEdit    bool
}

type pageHistoryAds struct {
	Rows []rowAds
}

type rowAds struct {
	ID        uint
	Text      string
	ImageID   string
	Status    uint8
	CreatedAt time.Time
}

func HandleMenu(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Добавить объявление", "AddAds"),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Мои объявления", "AdsHistory"),
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Назад", "StartMenu"),
		},
	)
	context.UpdateUserLevel(userID, ctx, 0)
	msg := tgbotapi.NewEditMessageTextAndMarkup(userID, state.MessageID, config.GlobalSettings.Texts.AddsMenu, inlineKeyboard)
	ctx.BotAPI.Send(msg)

}
func HandleSelectADSHistory(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 8)
	_, exist := state.Data["AdsHistory"]
	var rows [][]tgbotapi.InlineKeyboardButton
	if !exist {
		var ads []models.Advertisement
		db.DB.Joins(`JOIN "Users" ON "Users"."id" = "Advertisements"."user_id"`).
			Where(`"Users"."telegram_id" = ?`, userID).
			Order(`"Advertisements"."created_at" DESC`).
			Find(&ads)

		pageSize := 10
		state.Data["AdsHistoryPage"] = uint(0)
		state.Data["AdsHistory"] = make(map[uint]pageHistoryAds)
		pages := state.Data["AdsHistory"].(map[uint]pageHistoryAds)
		for page := 0; page < (len(ads)+pageSize-1)/pageSize; page++ {
			var _Ads []rowAds
			start := page * pageSize
			end := start + pageSize
			if end > len(ads) {
				end = len(ads)
			}
			for _, ad := range ads[start:end] {
				_Ads = append(_Ads, rowAds{ID: ad.ID, Text: ad.Text, Status: ad.Status, CreatedAt: ad.CreatedAt, ImageID: ad.ImageID})
			}

			pages[uint(page)] = pageHistoryAds{Rows: _Ads}

		}
		state.Data["AdsHistory"] = pages
	}
	pages := state.Data["AdsHistory"].(map[uint]pageHistoryAds)
	currentPage := state.Data["AdsHistoryPage"].(uint)
	page := pages[currentPage].Rows
	for i := 0; i < len(page); i += 2 {
		var titleI string = page[i].CreatedAt.UTC().Format("2006-01-02")
		switch page[i].Status {
		case 0:
			titleI += " ⏳"
		case 1:
			titleI += " ✅"
		case 2:
			titleI += " ❌"
		}
		if i+1 < len(page) {
			var titleI1 string = page[i+1].CreatedAt.UTC().Format("2006-01-02")
			switch page[i+1].Status {
			case 0:
				titleI1 += " ⏳"
			case 1:
				titleI1 += " ✅"
			case 2:
				titleI1 += " ❌"
			}
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(titleI, "Ad_"+strconv.Itoa(int(page[i].ID))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i))),
				tgbotapi.NewInlineKeyboardButtonData(titleI1, "Ad_"+strconv.Itoa(int(page[i+1].ID))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i+1))),
			))

		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(titleI, "Ad_"+strconv.Itoa(int(page[i].ID))+"_"+strconv.Itoa(int(currentPage))+"_"+strconv.Itoa(int(i))),
			))
		}
	}
	if len(pages)-1 > int(currentPage) && currentPage != 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("«", "backAds"), tgbotapi.NewInlineKeyboardButtonData("»", "nextAds")))
	} else if len(pages)-1 > int(currentPage) && currentPage == 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("»", "nextAds")))
	} else if len(pages)-1 == int(currentPage) && len(pages) != 1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("«", "backAds")))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		state.MessageID,
		"История объявлений",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}
func HandleViwerAdsHistory(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 9)
	pages := state.Data["AdsHistory"].(map[uint]pageHistoryAds)
	Indexes := strings.Split(update.CallbackQuery.Data, "_")
	pageIndex, _ := strconv.Atoi(Indexes[2])
	adsIndex, _ := strconv.Atoi(Indexes[3])
	ads := pages[uint(pageIndex)].Rows[adsIndex]
	var rows [][]tgbotapi.InlineKeyboardButton
	if ads.ImageID == "" {
		state.Data["MessageIdPhoto"] = 0
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			update.CallbackQuery.Message.Chat.ID,
			state.MessageID,
			ads.Text,
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		msg.ParseMode = "HTML"
		ctx.BotAPI.Send(msg)
		return
	}
	deleteMsg := tgbotapi.DeleteMessageConfig{
		ChatID:    userID,
		MessageID: state.MessageID,
	}
	ctx.BotAPI.Send(deleteMsg)
	photoConfig := tgbotapi.NewPhoto(userID, tgbotapi.FileID(ads.ImageID))
	message, _ := ctx.BotAPI.Send(photoConfig)
	state.Data["MessageIdPhoto"] = message.MessageID
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	msg := tgbotapi.NewMessage(
		update.CallbackQuery.Message.Chat.ID,
		ads.Text,
	)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ParseMode = "HTML"
	ctx.SendMessage(msg)
}
func HandleSelectADS(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	// var Ad models.Advertisement
	// var user models.User
	// db.DB.Where("telegram_id = ?", userID).First(&user)

	// threeHoursAgo := time.Now().Add(-3 * time.Hour)
	// db.DB.Model(&models.Advertisement{}).
	// 	Where("user_id = ? AND status IN (?) AND created_at >= ?", user.ID, []uint8{0, 1}, threeHoursAgo).
	// 	First(&Ad)
	// timeLimit := 3 * time.Hour
	// remainingTime := timeLimit - time.Since(Ad.CreatedAt)
	// if remainingTime > 0 {
	// 	hours := int(remainingTime.Hours())
	// 	minutes := int(remainingTime.Minutes()) % 60
	// 	message := fmt.Sprintf("Вы сможете создать новое объявление через %d часа %d минут.", hours, minutes)
	// 	alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, message)
	// 	alert.ShowAlert = false
	// 	ctx.BotAPI.Request(alert)
	// 	return

	// }
	context.UpdateUserLevel(userID, ctx, 1)

	var types []models.AdvertisementType
	db.DB.Find(&types)

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range types {
		if t.IsFree {
			button := tgbotapi.NewInlineKeyboardButtonData(t.Name, "newAds_"+strconv.Itoa(int(t.ID)))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
		} else {
			button := tgbotapi.NewInlineKeyboardButtonData(t.Name+" ("+strconv.Itoa(int(t.Cost))+"₩)", "newAds_"+strconv.Itoa(int(t.ID)))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
		}

	}
	button := tgbotapi.NewInlineKeyboardButtonData("Назад", "back")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		state.MessageID,
		"Выберете тип объявления",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}

func HandleAddAds(update *tgbotapi.Update, ctx *context.Context, typeID string) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)

	if state.Data["AdsInputs"] == nil {
		var inputs []models.AdvertisementInputs

		typeIDInt, _ := strconv.Atoi(typeID)
		typeIDUint := uint(typeIDInt)
		var Type models.AdvertisementType
		db.DB.Where(&models.AdvertisementType{ID: typeIDUint}).First(&Type)
		if !Type.IsFree {
			var User models.User
			db.DB.Where(&models.User{TelegramID: userID}).First(&User)
			if User.Balance < Type.Cost {
				message := "Недостаточно средств на балансе!"
				alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, message)
				alert.ShowAlert = false
				ctx.BotAPI.Request(alert)
				return
			}
		}
		db.DB.Where(&models.AdvertisementInputs{TypeID: typeIDUint}).Order("priority asc").Find(&inputs)

		resultMap := make(map[uint]AdsInputs)

		for _, input := range inputs {
			adsInput := AdsInputs{
				ID:       input.ID,
				Optional: input.Optional,
				Name:     input.Name,
				Activate: false,
				Value:    "",
				Priority: input.Priority,
			}
			resultMap[input.ID] = adsInput
		}
		state.Data["ActivType"] = typeIDUint
		state.Data["AdsInputs"] = resultMap
		state.Data["AdsPhoto"] = AdsPhoto{
			ID:        "",
			Activate:  false,
			MessageId: 0,
			ActivStep: 0,
		}
	}
	context.UpdateUserLevel(userID, ctx, 2)
	adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
	photo, _ := state.Data["AdsPhoto"].(AdsPhoto)
	var sortedInputs []AdsInputs
	for _, input := range adsInputs {
		sortedInputs = append(sortedInputs, input)
	}

	sort.Slice(sortedInputs, func(i, j int) bool {
		if sortedInputs[i].Priority == sortedInputs[j].Priority {
			return sortedInputs[i].ID < sortedInputs[j].ID
		}
		return sortedInputs[i].Priority < sortedInputs[j].Priority
	})

	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, input := range sortedInputs {
		callbackData := "AddInput_" + strconv.Itoa(int(input.ID))
		var suffux string
		if input.Activate {
			suffux = "✅"
		} else {
			if input.Optional {
				suffux = "❔"
			} else {
				suffux = "❌"
			}
		}
		button := tgbotapi.NewInlineKeyboardButtonData(input.Name+" • "+suffux, callbackData)

		row = append(row, button)

		if (uint(i)+1)%2 == 0 || uint(i) == uint(len(sortedInputs))-1 {
			rows = append(rows, row)
			row = nil
		}
	}
	var Types models.AdvertisementType
	db.DB.Where(&models.AdvertisementType{ID: state.Data["ActivType"].(uint)}).First(&Types)
	var messageText string = "<b>" + Types.Name + "</b>\n\n"
	for _, input := range sortedInputs {
		var Value string
		if input.Value == "" {
			if input.Optional {
				Value = "❔"
			} else {
				Value = "❌"
			}
		} else {
			Value = input.Value
		}
		messageText += "<b>" + input.Name + "</b>: " + Value + "\n\n"
	}

	if photo.Activate {
		suffux := "✅"
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Фото • "+suffux, "AddPhoto")))
	} else {
		suffux := "❔"
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Фото • "+suffux, "AddPhoto")))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пред просмотр", "preViwe")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сохранить", "Save")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		state.MessageID,
		messageText,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}
func HandleBackAds(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	text := "Вы уверены что хотите уйти? всё что вы ввели не сохранится"
	context.UpdateUserLevel(userID, ctx, 7)
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", "Delete")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		state.MessageID,
		text,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}
func HandleSaveAds(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 6)
	adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
	photo, _ := state.Data["AdsPhoto"].(AdsPhoto)
	typeID := state.Data["ActivType"].(uint)
	var rows [][]tgbotapi.InlineKeyboardButton
	var sortedInputs []AdsInputs
	for _, input := range adsInputs {
		sortedInputs = append(sortedInputs, input)
	}
	var Types models.AdvertisementType
	db.DB.Where(&models.AdvertisementType{ID: typeID}).First(&Types)

	sort.Slice(sortedInputs, func(i, j int) bool {
		if sortedInputs[i].Priority == sortedInputs[j].Priority {
			return sortedInputs[i].ID < sortedInputs[j].ID
		}
		return sortedInputs[i].Priority < sortedInputs[j].Priority
	})
	var messageText string = "<b>" + Types.Name + "</b>\n\n"
	for _, input := range sortedInputs {
		var Value string
		if input.Value == "" {
			if input.Optional {
				continue
			} else {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
				msg := tgbotapi.NewEditMessageTextAndMarkup(
					update.CallbackQuery.Message.Chat.ID,
					state.MessageID,
					"Вы не ввели все обязательные поля!",
					tgbotapi.NewInlineKeyboardMarkup(rows...),
				)
				ctx.BotAPI.Send(msg)
				return
			}
		} else {
			Value = input.Value
		}
		messageText += "<b>" + input.Name + "</b>: " + Value + "\n\n"
	}
	var User models.User

	db.DB.Where(&models.User{TelegramID: userID}).First(&User)
	var CostUser uint = 0
	if !Types.IsFree {
		CostUser = User.Balance - Types.Cost
		db.DB.Model(&models.User{}).
			Where("id = ?", uint(User.ID)).
			Updates(map[string]interface{}{
				"balance": User.Balance - Types.Cost,
			})
	}
	if photo.Activate {
		db.DB.Save(&models.Advertisement{UserID: uint(User.ID), Text: messageText, ImageID: photo.ID, Status: 0, CostUser: CostUser})
	} else {
		db.DB.Save(&models.Advertisement{UserID: uint(User.ID), Text: messageText, Status: 0, CostUser: CostUser})
	}
	delete(state.Data, "AdsInputs")
	delete(state.Data, "AdsPhoto")
	delete(state.Data, "ActivType")
	HandleMenu(update, ctx)

}

func HandlePreWive(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 5)
	adsInputs, _ := state.Data["AdsInputs"].(map[uint]AdsInputs)
	photo, _ := state.Data["AdsPhoto"].(AdsPhoto)
	typeID := state.Data["ActivType"].(uint)

	var sortedInputs []AdsInputs
	for _, input := range adsInputs {
		sortedInputs = append(sortedInputs, input)
	}

	sort.Slice(sortedInputs, func(i, j int) bool {
		if sortedInputs[i].Priority == sortedInputs[j].Priority {
			return sortedInputs[i].ID < sortedInputs[j].ID
		}
		return sortedInputs[i].Priority < sortedInputs[j].Priority
	})
	var Types models.AdvertisementType
	db.DB.Where(&models.AdvertisementType{ID: typeID}).First(&Types)
	var messageText string = "<b>" + Types.Name + "</b>\n\n"
	for _, input := range sortedInputs {
		var Value string
		if input.Value == "" {
			if input.Optional {
				continue
			} else {
				Value = "❌"
			}
		} else {
			Value = input.Value
		}
		messageText += "<b>" + input.Name + "</b>: " + Value + "\n\n"
	}
	if photo.Activate {
		photoConfig := tgbotapi.NewPhoto(userID, tgbotapi.FileID(photo.ID))
		message, _ := ctx.BotAPI.Send(photoConfig)
		photo.MessageId = message.MessageID
		state.Data["AdsPhoto"] = photo
	}
	deleteMsg1 := tgbotapi.DeleteMessageConfig{
		ChatID:    userID,
		MessageID: state.MessageID,
	}
	ctx.BotAPI.Send(deleteMsg1)
	msg := tgbotapi.NewMessage(
		userID,
		messageText,
	)
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Назад", "back")))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ParseMode = "HTML"
	ctx.SendMessage(msg)
}

func HandleAddPhoto(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	var PhotoID string = ""
	if update.Message != nil {
		userID = update.Message.Chat.ID
		if len(update.Message.Photo) > 0 {
			largestPhoto := update.Message.Photo[len(update.Message.Photo)-1]
			PhotoID = largestPhoto.FileID
		}
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg1)

	} else {
		userID = update.CallbackQuery.From.ID
	}
	state := context.GetUserState(userID, ctx)
	photo, _ := state.Data["AdsPhoto"].(AdsPhoto)
	context.UpdateUserLevel(userID, ctx, 4)
	var rows [][]tgbotapi.InlineKeyboardButton
	if photo.Activate && photo.IsEdit {
		switch photo.ActivStep {
		case 0:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "AddPhoto" {
					photoConfig := tgbotapi.NewPhoto(userID, tgbotapi.FileID(photo.ID))
					message, _ := ctx.BotAPI.Send(photoConfig)
					photo.MessageId = message.MessageID
					photo.ActivStep = 1
					state.Data["AdsPhoto"] = photo
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать", "Edit")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", "Delete")))
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "back")))
					msg := tgbotapi.NewMessage(
						userID,
						"❔ Вы хотите изменить фотографию?",
					)
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
					deleteMsg1 := tgbotapi.DeleteMessageConfig{
						ChatID:    userID,
						MessageID: state.MessageID,
					}
					msg.ParseMode = "HTML"
					ctx.BotAPI.Send(deleteMsg1)
					ctx.SendMessage(msg)
					return
				}
			}
		case 1:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Edit" {
					photo.ActivStep = 0
					photo.IsEdit = false
					state.Data["AdsPhoto"] = photo
				}
				if update.CallbackQuery.Data == "Delete" {
					photo.ActivStep = 2
					state.Data["AdsPhoto"] = photo
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
				}
			}
		case 2:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Delete" {
					photo.ActivStep = 0
					photo.IsEdit = false
					photo.Activate = false
					state.Data["AdsPhoto"] = photo
					HandleAddAds(update, ctx, "0")
					return
				}
			}

		}
	}
	if !photo.IsEdit || !photo.Activate && photo.IsEdit {
		switch photo.ActivStep {
		case 0:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "AddPhoto" || update.CallbackQuery.Data == "Edit" {
					photo.ActivStep = 1
					state.Data["AdsPhoto"] = photo
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
					msg := tgbotapi.NewEditMessageTextAndMarkup(
						update.CallbackQuery.Message.Chat.ID,
						state.MessageID,
						"❔Отправьте фотографию",
						tgbotapi.NewInlineKeyboardMarkup(rows...),
					)
					ctx.BotAPI.Send(msg)
					return
				}
			}
		case 1:
			if update.Message != nil && PhotoID != "" {
				photo.ActivStep = 2
				photoConfig := tgbotapi.NewPhoto(userID, tgbotapi.FileID(PhotoID))
				message, _ := ctx.BotAPI.Send(photoConfig)

				photo.MessageId = message.MessageID
				photo.IDpre = PhotoID
				state.Data["AdsPhoto"] = photo
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📋 Сохранить", "Save")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🚫 Отмена ", "back")))
				msg := tgbotapi.NewMessage(
					userID,
					"📋 Сохранить фотографию?",
				)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
				deleteMsg1 := tgbotapi.DeleteMessageConfig{
					ChatID:    userID,
					MessageID: state.MessageID,
				}
				ctx.BotAPI.Send(deleteMsg1)
				ctx.SendMessage(msg)
				return
			}
		case 2:
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "Save" {
					photo.Activate = true
					photo.IsEdit = true
					photo.ActivStep = 0
					photo.ID = photo.IDpre
					if photo.MessageId != 0 {
						deleteMsg1 := tgbotapi.DeleteMessageConfig{
							ChatID:    userID,
							MessageID: photo.MessageId,
						}
						photo.MessageId = 0
						ctx.BotAPI.Send(deleteMsg1)
					}
					state.Data["AdsPhoto"] = photo
					HandleAddAds(update, ctx, "0")
					return
				}
			}

		}
	}

}

// func HandleGetValue(update *tgbotapi.Update, ctx *context.Context) {
// 	userID := update.Message.From.ID
// 	deleteMsg1 := tgbotapi.DeleteMessageConfig{
// 		ChatID:    userID,
// 		MessageID: update.Message.MessageID,
// 	}
// 	ctx.BotAPI.Send(deleteMsg1)
// 	value := update.Message.Text
// 	entities := update.Message.Entities
// 	state := context.GetUserState(userID, ctx)

// 	// inputID := state.Data["ActiveInput"]
// 	// adsInputs := state.Data["AdsInputs"].(map[uint]AdsInputs)
// 	formatetText := utilits.ApplyFormatting(value, entities)

// 	var rows [][]tgbotapi.InlineKeyboardButton
// 	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Сохранить", "back")))
// 	msg := tgbotapi.NewEditMessageTextAndMarkup(
// 		userID,
// 		state.MessageID,
// 		formatetText,
// 		tgbotapi.NewInlineKeyboardMarkup(rows...),
// 	)
// 	msg.ParseMode = "HTML"
// 	ctx.BotAPI.Send(msg)
// }
