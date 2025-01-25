package globalsetingsroutes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/utilits"

	"github.com/go-chi/chi/v5"
)

func Docs(r chi.Router) {
	r.Get("/Docs", func(w http.ResponseWriter, r *http.Request) {
		videoUrl, err := utilits.GetPhotoLink(config.GlobalSettings.Docs.VideoID)
		if err != nil {
			fmt.Println(err)
		}
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]string{
				"Video": videoUrl,
				"Text":  config.GlobalSettings.Docs.Text,
			},
		})
	})
	r.Put("/Docs", func(w http.ResponseWriter, r *http.Request) {
		const MaxUploadSize = 100 << 20
		r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

		if err := r.ParseMultipartForm(MaxUploadSize); err != nil && err != http.ErrNotMultipart {
			http.Error(w, "Ошибка обработки данных: "+err.Error(), http.StatusBadRequest)
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
				http.Error(w, "Ошибка обработки видео: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		setings := config.GlobalSettings
		if fileID != "" {
			setings.Docs.VideoID = fileID
		}
		setings.Docs.Text = r.FormValue("text")
		config.Save(setings)

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Данные успешно обновлены",
		})
	})

}
