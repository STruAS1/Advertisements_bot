package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/utilits"

	"github.com/go-chi/chi/v5"
)

type DocsRequest struct {
	Video string `json:"Video"`
	Text  string `json:"Text"`
}

func Docs(r chi.Router) {
	r.Get("/Docs", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]string{
				"Video": config.GlobalSettings.Docs.VideoUrl,
				"Text":  config.GlobalSettings.Docs.Text,
			},
		})
	})
	r.Put("/Docs", func(w http.ResponseWriter, r *http.Request) {
		var Docs DocsRequest
		if err := json.NewDecoder(r.Body).Decode(&Docs); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		if setings.Docs.VideoUrl != Docs.Video {
			setings.Docs.VideoUrl = Docs.Video
			FileId, _ := utilits.FetchVideoToMemory(Docs.Video)
			setings.Docs.VideoID = FileId

		}
		setings.Docs.Text = Docs.Text
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}
