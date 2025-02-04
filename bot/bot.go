package bot

import (
	"log"
	"tgbotBARAHOLKA/bot/context"
	"tgbotBARAHOLKA/bot/handlers"
	"tgbotBARAHOLKA/config"

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
				log.Println(update.Message.Chat.ID)
				log.Println(update.Message.SenderChat.UserName)
				if update.Message.Chat.ID == cfg.Bot.CommentChatId && update.Message.SenderChat.UserName == cfg.Bot.ChannelId {
					config.LastUpdateFromChannel = &update
					continue
				}
			}
		}
		if (update.Message != nil && update.Message.Chat.Type == "private") || (update.CallbackQuery != nil && update.CallbackQuery.Message.Chat.Type == "private") {
			handlers.HandleUpdate(&update, ctx)
		}
	}
}
