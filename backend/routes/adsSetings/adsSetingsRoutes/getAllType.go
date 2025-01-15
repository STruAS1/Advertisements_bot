package adsSetingsRoutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	"github.com/go-chi/chi/v5"
)

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

func GetAllAdvertisementTypes(r chi.Router) {
	r.Get("/Types", func(w http.ResponseWriter, r *http.Request) {
		var types []models.AdvertisementType
		if err := db.DB.Find(&types).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch data from database",
			})
			return
		}

		if len(types) == 0 {
			writeJSON(w, http.StatusOK, SuccessResponse{
				Message: "Ok",
				Data:    []models.AdvertisementType{},
			})
			return
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data:    types,
		})
	})
}
