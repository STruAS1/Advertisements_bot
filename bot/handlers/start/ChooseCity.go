package start

import (
	"fmt"
	"strconv"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CitiesRow struct {
	Id       uint
	Title    string
	IsActive bool
}

type CitiesPage struct {
	Cities []CitiesRow
}

type ActiveChooseCityType struct {
	CitiesPages map[uint]CitiesPage
	CurentPage  uint
	ActiveStep  uint
}

func ChooseCityHandler(update *tgbotapi.Update, ctx *context.Context) {
	var userID int64
	var value string
	if update.Message != nil {
		userID = update.Message.Chat.ID
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: update.Message.MessageID,
		}
		value = update.Message.Text
		ctx.BotAPI.Send(deleteMsg1)

	} else {
		userID = update.CallbackQuery.From.ID
	}
	context.UpdateUserLevel(userID, ctx, 4)
	state := context.GetUserState(userID, ctx)
	_, exist := state.Data["ActiveChooseCity"]
	if !exist {
		data := ActiveChooseCityType{}
		data.CitiesPages = make(map[uint]CitiesPage)
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

			data.CitiesPages[uint(page)] = CitiesPage{Cities: _Cities}

		}

		state.Data["ActiveChooseCity"] = data
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	ActiveChooseCity, _ := state.Data["ActiveChooseCity"].(ActiveChooseCityType)
	if update.CallbackQuery != nil && update.CallbackQuery.Data == "BackToList" {
		ActiveChooseCity.ActiveStep = 0
		state.Data["ActiveChooseCity"] = ActiveChooseCity
	}
	switch ActiveChooseCity.ActiveStep {
	case 0:
		CallbackQuery := strings.Split(update.CallbackQuery.Data, "_")
		if CallbackQuery[0] == "ChooseCity" || update.CallbackQuery.Data == "Edit" || update.CallbackQuery.Data == "nextCity" || update.CallbackQuery.Data == "backCity" || update.CallbackQuery.Data == "search" || update.CallbackQuery.Data == "BackToList" || CallbackQuery[0] == "City" || update.CallbackQuery.Data == "Save" || update.CallbackQuery.Data == "menuCityInfo" {
			if update.CallbackQuery.Data == "nextCity" && len(ActiveChooseCity.CitiesPages)-1 != int(ActiveChooseCity.CurentPage) {
				ActiveChooseCity.CurentPage++
				state.Data["ActiveChooseCity"] = ActiveChooseCity
			}
			if update.CallbackQuery.Data == "backCity" && int(ActiveChooseCity.CurentPage) != 0 {
				ActiveChooseCity.CurentPage--
				state.Data["ActiveChooseCity"] = ActiveChooseCity
			}
			if update.CallbackQuery.Data == "menuCityInfo" {
				alert := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Выберете город из списка")
				alert.ShowAlert = false
				ctx.BotAPI.Request(alert)
				return
			}
			ActiveChooseCity.ActiveStep = 0
			if update.CallbackQuery.Data == "search" {
				var textActiveCities string = ""
				ActiveChooseCity.ActiveStep = 1
				state.Data["ActiveChooseCity"] = ActiveChooseCity
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
			currentPage := ActiveChooseCity.CurentPage
			page := ActiveChooseCity.CitiesPages[currentPage].Cities
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
			if len(ActiveChooseCity.CitiesPages)-1 > int(currentPage) && currentPage != 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "backCity"), tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search"), tgbotapi.NewInlineKeyboardButtonData("Дальше »", "nextCity")))
			} else if len(ActiveChooseCity.CitiesPages)-1 > int(currentPage) && currentPage == 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search"), tgbotapi.NewInlineKeyboardButtonData("Дальше »", "nextCity")))
			} else if len(ActiveChooseCity.CitiesPages)-1 == int(currentPage) && len(ActiveChooseCity.CitiesPages) != 1 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад", "backCity"), tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search")))
			} else {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔎 Поиск", "search")))
			}
			text := "🏙️ Выберите город из списка!" + "\n"
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

			for pageID, page := range ActiveChooseCity.CitiesPages {
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

			if len(sortedCities) == 0 {
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Совпадений не найдено", "NoResult")))
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("« Назад к списку", "BackToList")))
				text := "Введите название Города" + "\n"
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
				text := "Введите название Города" + "\n"
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
