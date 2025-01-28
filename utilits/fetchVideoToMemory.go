package utilits

import (
	"bytes"
	"errors"
	"tgbotBARAHOLKA/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SaveAndSendVideoToTelegram(fileName string, fileData []byte) (string, error) {
	cfg := config.LoadConfig()
	botAPI, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return "", errors.New("ошибка подключения к Telegram API: " + err.Error())
	}

	video := tgbotapi.NewVideo(1062226084, tgbotapi.FileReader{
		Name:   fileName,
		Reader: bytes.NewReader(fileData),
	})
	video.ParseMode = "HTML"

	videoMessage, err := botAPI.Send(video)
	if err != nil {
		return "", errors.New("ошибка отправки видео в Telegram: " + err.Error())
	}

	return videoMessage.Video.FileID, nil
}
