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

		handlers.HandleUpdate(&update, ctx)
	}
}
