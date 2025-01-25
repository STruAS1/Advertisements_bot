package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type ButtonRquest struct {
	ID   uint16 `json:"ID"`
	Text string `json:"ButtonText"`
}

func Buttons(r chi.Router) {
	r.Get("/Buttons", func(w http.ResponseWriter, r *http.Request) {
		ButtonsResponce := make([]map[string]interface{}, len(config.GlobalSettings.Buttons))
		for i, Button := range config.GlobalSettings.Buttons {
			ButtonsResponce[i] = map[string]interface{}{
				"ID":          i,
				"ButtonText":  Button.ButtonText,
				"Description": Button.Discription,
			}
		}
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    ButtonsResponce,
		})
	})
	r.Put("/Buttons", func(w http.ResponseWriter, r *http.Request) {
		var ButtonRquest ButtonRquest
		if err := json.NewDecoder(r.Body).Decode(&ButtonRquest); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if len(config.GlobalSettings.Buttons) <= int(ButtonRquest.ID) {
			http.Error(w, "Index out of range", http.StatusBadRequest)
			return
		}
		settings := config.GlobalSettings
		settings.Buttons[ButtonRquest.ID].ButtonText = ButtonRquest.Text
		config.Save(settings)
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
		})
	})

}
