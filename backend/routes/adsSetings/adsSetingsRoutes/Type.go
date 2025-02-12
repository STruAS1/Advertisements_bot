package adsSetingsRoutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	"github.com/go-chi/chi/v5"
)

type TypeRequest struct {
	IsFree                    bool   `json:"IsFree"`
	Cost                      uint   `json:"Cost"`
	Name                      string `json:"Name"`
	Priority                  uint   `json:"priority"`
	AutoPost                  bool   `json:"AutoPost"`
	HasLimit                  bool   `json:"HasLimit"`
	DayLimit                  uint   `json:"DayLimit"`
	DayLimitWithVerification  uint   `json:"DayLimitWithVerification"`
	HourLimit                 uint   `json:"HourLimit"`
	HourLimitWithVerification uint   `json:"HourLimitWithVerification"`
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

func Type(r chi.Router) {
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

	r.Put("/Type", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateData TypeRequest
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := db.DB.Model(&models.AdvertisementType{}).
			Where("id = ?", uint(id)).
			Updates(map[string]interface{}{
				"is_free":                      updateData.IsFree,
				"name":                         updateData.Name,
				"Cost":                         updateData.Cost,
				"priority":                     updateData.Priority,
				"auto_post":                    updateData.AutoPost,
				"has_limit":                    updateData.HasLimit,
				"day_limit":                    updateData.DayLimit,
				"day_limit_with_verification":  updateData.DayLimitWithVerification,
				"hour_limit":                   updateData.HourLimit,
				"hour_limit_with_verification": updateData.HourLimitWithVerification,
			}).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to update record",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
	r.Get("/Type", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}

		var advType models.AdvertisementType
		if err := db.DB.Where(models.AdvertisementType{ID: uint(id)}).First(&advType).Error; err != nil {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Message: "Record not found",
			})
			return
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "oK",
			Data:    advType,
		})
	})

	r.Post("/Type", func(w http.ResponseWriter, r *http.Request) {
		var createData TypeRequest
		if err := json.NewDecoder(r.Body).Decode(&createData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if createData.Name == "" {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Название не может быть пустым ",
			})
		}
		newType := models.AdvertisementType{
			IsFree:                    createData.IsFree,
			Name:                      createData.Name,
			Cost:                      createData.Cost,
			Priority:                  createData.Priority,
			AutoPost:                  createData.AutoPost,
			HasLimit:                  createData.HasLimit,
			DayLimit:                  createData.DayLimit,
			DayLimitWithVerification:  createData.DayLimitWithVerification,
			HourLimit:                 createData.HourLimit,
			HourLimitWithVerification: createData.HourLimitWithVerification,
		}

		if err := db.DB.Create(&newType).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to create new type",
			})
			return
		}

		writeJSON(w, http.StatusCreated, SuccessResponse{
			Message: "OK",
			Data:    newType,
		})
	})

	r.Delete("/Type", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		if err := db.DB.Where(models.AdvertisementType{ID: uint(id)}).Delete(&models.AdvertisementType{}).Error; err != nil {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Message: "Record not found",
			})
			return
		}
		writeJSON(w, http.StatusCreated, SuccessResponse{
			Message: "OK",
			Data:    nil,
		})
	})
}
