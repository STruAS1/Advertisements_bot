package adsSetingsRoutes

import (
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	"github.com/go-chi/chi/v5"
)

func GetAllAdvertisementInputs(r chi.Router) {
	r.Get("/Inputs", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		TypeIDSTring := queryParams.Get("ID")
		TypeID, err := strconv.ParseUint(TypeIDSTring, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var Inputs []models.AdvertisementInputs
		if err := db.DB.Where(models.AdvertisementInputs{TypeID: uint(TypeID)}).Find(&Inputs).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch data from database",
			})
			return
		}

		if len(Inputs) == 0 {
			writeJSON(w, http.StatusOK, SuccessResponse{
				Message: "Ok",
				Data:    []models.AdvertisementType{},
			})
			return
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data:    Inputs,
		})
	})
}
