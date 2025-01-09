package start

// import (
// 	"tgbotBARAHOLKA/bot/context"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func HandleMenu(update *tgbotapi.Update, ctx *context.Context) {
// 	userID := update.CallbackQuery.From.ID
// 	state := context.GetUserState(userID, ctx)
// 	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		[]tgbotapi.InlineKeyboardButton{
// 			tgbotapi.NewInlineKeyboardButtonData("Подписаться на канал", "Add"),
// 		},
// 		// []tgbotapi.InlineKeyboardButton{
// 		// 	tgbotapi.NewInlineKeyboardButtonData("Подписался", "cehk_sub"),
// 		// },
// 	)
// 	msg := tgbotapi.NewEditMessageTextAndMarkup(userID, state.MessageID, "Меню", inlineKeyboard)
// 	ctx.BotAPI.Send(msg)

// }
