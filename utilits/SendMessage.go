package utilits

import (
	"fmt"
	"tgbotBARAHOLKA/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageToChnale(message, photoUrl string) int {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)
	var ms tgbotapi.Message
	var err error
	if photoUrl != "" {
		msg := tgbotapi.NewPhotoToChannel(cfg.Bot.ChannelId, tgbotapi.FileURL(photoUrl))
		msg.Caption = message
		msg.ParseMode = "HTML"
		ms, err = botAPI.Send(msg)

	} else {
		msg := tgbotapi.NewMessageToChannel(cfg.Bot.ChannelId, message)
		msg.ParseMode = "HTML"
		ms, err = botAPI.Send(msg)
	}
	if err != nil {
		fmt.Println(err)
	}
	return ms.MessageID
}

func SendMessageToUser(message string, userID int64) int {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)
	var ms tgbotapi.Message
	var err error

	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = "HTML"
	ms, err = botAPI.Send(msg)

	if err != nil {
		fmt.Println(err)
	}
	return ms.MessageID
}

func DeleteMessageFromChanel(massgeID int) error {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)

	deleteMsg1 := tgbotapi.DeleteMessageConfig{
		ChannelUsername: cfg.Bot.ChannelId,
		MessageID:       massgeID,
	}
	_, err := botAPI.Send(deleteMsg1)
	if err != nil {
		return err
	}
	return nil
}
