package utilits

import (
	"fmt"
	"log"
	"tgbotBARAHOLKA/config"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageToChnale(message, photoUrl string) (int, int) {
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
		msg.DisableWebPagePreview = true
		ms, err = botAPI.Send(msg)
	}
	if err != nil {
		fmt.Println(err)
	}
	var commentID int
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		if config.LastUpdateFromChannel != nil && config.LastUpdateFromChannel.Message != nil && config.LastUpdateFromChannel.Message.Text == RemoveHTMLTags(message) {
			commentID = config.LastUpdateFromChannel.Message.MessageID
			config.LastUpdateFromChannel = nil
			break
		}
	}
	return ms.MessageID, commentID
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

func CheckAndKickUserFromChannel(userID int64) error {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)

	chatMember, err := botAPI.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			SuperGroupUsername: cfg.Bot.ChannelId,
			UserID:             userID,
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка получения информации о пользователе: %v", err)
	}

	if chatMember.Status == "left" || chatMember.Status == "kicked" {
		kickConfig := tgbotapi.KickChatMemberConfig{
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				SuperGroupUsername: cfg.Bot.ChannelId,
				UserID:             userID,
			},
		}
		if _, err := botAPI.Send(kickConfig); err != nil {
			return fmt.Errorf("ошибка кика пользователя: %v", err)
		}
		return nil
	}

	return fmt.Errorf("пользователь уже состоит в канале")
}

func UnbanUser(userID int64) error {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)
	unban := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			SuperGroupUsername: cfg.Bot.ChannelId,
			UserID:             userID,
		},
		OnlyIfBanned: true,
	}

	_, err := botAPI.Request(unban)
	if err != nil {
		log.Printf("Ошибка при разбане пользователя %d: %v", userID, err)
		return err
	}

	return nil
}
