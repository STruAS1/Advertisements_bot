package start

import (
	"fmt"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSelectDocs(update *tgbotapi.Update, ctx *context.Context) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 3)
	var rows [][]tgbotapi.InlineKeyboardButton
	docs := config.GlobalSettings.Docs
	for i, doc := range docs {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(doc.ButtonName, fmt.Sprintf("doc_%d", i))))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[5].ButtonText, "back")))
	msg := tgbotapi.NewEditMessageTextAndMarkup(
		userID,
		state.MessageID,
		config.GlobalSettings.Texts.DocsText,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)
	msg.ParseMode = "HTML"
	ctx.BotAPI.Send(msg)
}
func HandleDocs(update *tgbotapi.Update, ctx *context.Context, IndexDocs int) {
	userID := update.CallbackQuery.From.ID
	state := context.GetUserState(userID, ctx)
	context.UpdateUserLevel(userID, ctx, 5)
	if len(config.GlobalSettings.Docs) < IndexDocs+1 {
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	docs := config.GlobalSettings.Docs[IndexDocs]
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(config.GlobalSettings.Buttons[5].ButtonText, "back")))
	if docs.VideoID != "" {
		deleteMsg1 := tgbotapi.DeleteMessageConfig{
			ChatID:    userID,
			MessageID: state.MessageID,
		}
		ctx.BotAPI.Send(deleteMsg1)
		video := tgbotapi.NewVideo(userID, tgbotapi.FileID(docs.VideoID))
		video.ParseMode = "HTML"
		vidoeMassge, err := ctx.BotAPI.Send(video)
		if err != nil {
			fmt.Println(err)
		}
		state.Data["LastVideoMassgeID"] = vidoeMassge.MessageID
		msg := tgbotapi.NewMessage(userID, docs.Text+"ㅤ")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ParseMode = "HTML"
		ctx.SendMessage(msg)
	} else {
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			userID,
			state.MessageID,
			docs.Text+"ㅤ",
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		msg.ParseMode = "HTML"
		ctx.BotAPI.Send(msg)
	}
}
