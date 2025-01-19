package utilits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tgbotBARAHOLKA/config"
)

type GetFileResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID   string `json:"file_id"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}

func GetPhotoLink(fileID string) (string, error) {
	cfg := config.LoadConfig()
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", cfg.Bot.Token, fileID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("ошибка при запросе к Telegram API: %w", err)
	}
	defer resp.Body.Close()

	var getFileResponse GetFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&getFileResponse); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if !getFileResponse.Ok {
		return "", fmt.Errorf("не удалось получить файл: файл с id %s не найден", fileID)
	}

	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", cfg.Bot.Token, getFileResponse.Result.FilePath)
	return downloadURL, nil
}
