package utilits

import (
	"bytes"
	"io"
	"net/http"
	"tgbotBARAHOLKA/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FetchVideoToMemory(url string) (string, error) {
	cfg := config.LoadConfig()
	botAPI, _ := tgbotapi.NewBotAPI(cfg.Bot.Token)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", io.ErrUnexpectedEOF
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}

	fileName := "video.mp4"
	video := tgbotapi.NewVideo(1062226084, tgbotapi.FileReader{
		Name:   fileName,
		Reader: bytes.NewReader(buf.Bytes()),
	})
	video.ParseMode = "HTML"
	vidoeMassge, _ := botAPI.Send(video)
	return vidoeMassge.Video.FileID, nil
}
