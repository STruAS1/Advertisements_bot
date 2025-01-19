package adsSetingsRoutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type EditSufixRequest struct {
	Sufix string `json:"Sufix"`
}

func Sufix(r chi.Router) {
	r.Put("/Sufix", func(w http.ResponseWriter, r *http.Request) {
		var EditSufixData EditSufixRequest
		if err := json.NewDecoder(r.Body).Decode(&EditSufixData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Ads.Sufix = EditSufixData.Sufix
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})

	r.Get("/Sufix", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]string{
				"sufix": config.GlobalSettings.Ads.Sufix,
			},
		})
	})
}
