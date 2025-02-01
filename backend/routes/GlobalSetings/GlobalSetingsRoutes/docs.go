package globalsetingsroutes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/utilits"

	"github.com/go-chi/chi/v5"
)

func Docs(r chi.Router) {
	r.Get("/Docs", func(w http.ResponseWriter, r *http.Request) {
		DocsIndex := make([]map[string]interface{}, len(config.GlobalSettings.Docs))
		for i, Doc := range config.GlobalSettings.Docs {
			DocsIndex[i] = map[string]interface{}{
				"index":      i,
				"text":       Doc.ButtonName,
				"ButtonName": Doc.Text,
			}
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    DocsIndex,
		})
	})
	r.Put("/Docs", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("Index")
		Index, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		if len(config.GlobalSettings.Docs) <= int(Index) {
			http.Error(w, "Index out of range", http.StatusBadRequest)
			return
		}

		const MaxUploadSize = 100 << 20
		r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

		if err := r.ParseMultipartForm(MaxUploadSize); err != nil && err != http.ErrNotMultipart {
			http.Error(w, "Ошибка обработки данных: "+err.Error(), http.StatusBadRequest)
			return
		}

		newText := r.FormValue("text")
		newButtonName := r.FormValue("ButtonName")

		var fileID string
		if file, header, err := r.FormFile("video"); err == nil {
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
				return
			}

			fileID, err = utilits.SaveAndSendVideoToTelegram(header.Filename, buf.Bytes())
			if err != nil {
				http.Error(w, "Ошибка обработки видео: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		setings := config.GlobalSettings
		if fileID != "" {
			setings.Docs[Index].VideoID = fileID
		}
		if newText != "" {
			setings.Docs[Index].Text = newText
		}
		if newButtonName != "" {
			setings.Docs[Index].ButtonName = newButtonName
		}

		config.Save(setings)

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Данные успешно обновлены",
		})
	})
	r.Post("/Docs", func(w http.ResponseWriter, r *http.Request) {
		const MaxUploadSize = 100 << 20
		r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

		if err := r.ParseMultipartForm(MaxUploadSize); err != nil && err != http.ErrNotMultipart {
			http.Error(w, "Ошибка обработки данных: "+err.Error(), http.StatusBadRequest)
			return
		}

		ButtonName := r.FormValue("ButtonName")
		Text := r.FormValue("Text")
		if ButtonName == "" || Text == "" {
			http.Error(w, "ButtonName и Text обязательны", http.StatusBadRequest)
			return
		}

		var fileID string
		if file, header, err := r.FormFile("video"); err == nil {
			defer file.Close()

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
				return
			}

			fileID, err = utilits.SaveAndSendVideoToTelegram(header.Filename, buf.Bytes())
			if err != nil {
				http.Error(w, "Ошибка загрузки видео: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		newDoc := config.Docs{
			ButtonName: ButtonName,
			VideoID:    fileID,
			Text:       Text,
		}
		setings := config.GlobalSettings
		setings.Docs = append(setings.Docs, newDoc)
		config.Save(setings)

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Обучение успешно добавлено",
		})
	})
	r.Delete("/Docs", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("Index")

		Index, err := strconv.Atoi(idStr)
		if err != nil || Index < 0 {
			http.Error(w, "Invalid Index: must be a non-negative integer", http.StatusBadRequest)
			return
		}

		setings := config.GlobalSettings
		if Index >= len(setings.Docs) {
			http.Error(w, "Index out of range", http.StatusBadRequest)
			return
		}

		setings.Docs = append(setings.Docs[:Index], setings.Docs[Index+1:]...)
		config.Save(setings)

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Обучение успешно удалено",
		})
	})
	r.Put("/DocsText", func(w http.ResponseWriter, r *http.Request) {
		var EditTextData EditTextRequest
		if err := json.NewDecoder(r.Body).Decode(&EditTextData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Texts.DocsText = EditTextData.Text
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Get("/DocsText", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]string{
				"Text": config.GlobalSettings.Texts.DocsText,
			},
		})
	})
}
