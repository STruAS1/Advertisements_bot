package adsSetingsRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	"github.com/go-chi/chi/v5"
)

type UpdateInputRequest struct {
	Priority uint   `json:"priority"`
	Name     string `json:"name"`
	Options  string `json:"options"`
	Optional bool   `json:"optional"`
	InputID  uint   `json:"input_id"`
}

func UpdateAdvertisementInput(r chi.Router) {
	r.Put("/Input", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateData UpdateInputRequest
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := db.DB.Model(&models.AdvertisementInputs{}).
			Where("id = ?", uint(id)).
			Updates(map[string]interface{}{
				"priority": updateData.Priority,
				"name":     updateData.Name,
				"options":  updateData.Options,
				"optional": updateData.Optional,
				"input_id": updateData.InputID,
			}).Error; err != nil {
			http.Error(w, "Failed to update record", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}

func GetAdvertisementInputInfo(r chi.Router) {
	r.Get("/Input", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		IDSTring := queryParams.Get("ID")
		ID, err := strconv.ParseUint(IDSTring, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var Input models.AdvertisementInputs
		if err := db.DB.Where(models.AdvertisementInputs{ID: uint(ID)}).First(&Input).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch data from database",
			})
			return
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data:    Input,
		})
	})
}

type CreateInputRequest struct {
	Priority uint   `json:"priority"`
	Name     string `json:"name"`
	Options  string `json:"options"`
	Optional bool   `json:"optional"`
	InputID  uint   `json:"input_id"`
	TypeID   uint   `json:"type_id"`
}

func CreateAdvertisementInput(r chi.Router) {
	r.Post("/Input", func(w http.ResponseWriter, r *http.Request) {
		var createData CreateInputRequest
		if err := json.NewDecoder(r.Body).Decode(&createData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		newInput := models.AdvertisementInputs{
			Priority: createData.Priority,
			Name:     createData.Name,
			Options:  createData.Options,
			Optional: createData.Optional,
			InputID:  createData.InputID,
			TypeID:   createData.TypeID,
		}

		if err := db.DB.Create(&newInput).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to create new input",
			})
			return
		}

		writeJSON(w, http.StatusCreated, SuccessResponse{
			Message: "Ok",
			Data:    newInput,
		})
	})
}
