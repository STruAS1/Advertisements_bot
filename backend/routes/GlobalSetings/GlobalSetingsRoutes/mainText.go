package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type EditTextRequest struct {
	Text string `json:"Text"`
}
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type ErrorResponse struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func MainText(r chi.Router) {
	r.Put("/MainText", func(w http.ResponseWriter, r *http.Request) {
		var EditMainTextData EditTextRequest
		if err := json.NewDecoder(r.Body).Decode(&EditMainTextData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Texts.MainText = EditMainTextData.Text
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Get("/MainText", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]string{
				"Text": config.GlobalSettings.Texts.MainText,
			},
		})
	})
}

func AdsText(r chi.Router) {
	r.Put("/AdsText", func(w http.ResponseWriter, r *http.Request) {
		var EditTextData EditTextRequest
		if err := json.NewDecoder(r.Body).Decode(&EditTextData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Texts.AddsMenu = EditTextData.Text
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Get("/AdsText", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]string{
				"Text": config.GlobalSettings.Texts.AddsMenu,
			},
		})
	})
}
