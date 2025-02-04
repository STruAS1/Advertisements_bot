package bot

import (
	"fmt"
	"log"
	"strings"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(cfg *config.Config) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	ctx := context.NewContext(botAPI, cfg)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.Message != nil && update.Message.SenderChat != nil {

			if update.Message != nil && update.Message.SenderChat != nil {
				if update.Message.Chat.ID == cfg.Bot.CommentChatId && update.Message.SenderChat.UserName == strings.TrimPrefix(cfg.Bot.ChannelId, "@") {
					config.LastUpdateFromChannel = &update
					continue
				}
			}
		}
		if update.Message != nil && update.Message.Chat.ID == cfg.Bot.CommentChatId && update.Message.Text != "" {
			if update.Message.ReplyToMessage != nil {
				originalMessageID := update.Message.ReplyToMessage.MessageID
				fmt.Print(originalMessageID)
				var Ad models.Advertisement
				result := db.DB.Preload("User").Where(&models.Advertisement{CommentMsgId: originalMessageID}).First(&Ad)
				if result.Error != nil {
					msg := tgbotapi.NewMessage(int64(Ad.User.TelegramID), fmt.Sprintf("❗Новый комментарий:\n%s\n\n<a href='https://t.me/\u200B%s/%d>Объявление</a>", update.Message.Text, cfg.Bot.ChannelId, Ad.MassgeID))
					if _, err := botAPI.Send(msg); err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		if (update.Message != nil && update.Message.Chat.Type == "private") || (update.CallbackQuery != nil && update.CallbackQuery.Message.Chat.Type == "private") {
			handlers.HandleUpdate(&update, ctx)
		}
	}
}
