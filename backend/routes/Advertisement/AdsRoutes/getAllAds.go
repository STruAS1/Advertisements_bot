package adsroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"
	"tgbotBARAHOLKA/utilits"

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

func GetAllAdvertisements(r chi.Router) {
	r.Get("/ads", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		statusStr := queryParams.Get("status")
		limitStr := queryParams.Get("limit")
		pageStr := queryParams.Get("page")

		var status uint8
		var limit, page int
		var err error

		if statusStr != "" {
			statusUint64, err := strconv.ParseUint(statusStr, 10, 8)
			if err != nil || statusUint64 > 3 {
				http.Error(w, "Invalid status: must be between 0 and 3", http.StatusBadRequest)
				return
			}

			status = uint8(statusUint64)
		}

		if limitStr == "" {
			limit = 20
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				http.Error(w, "Invalid limit: must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		if pageStr == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page <= 0 {
				http.Error(w, "Invalid page: must be a positive integer", http.StatusBadRequest)
				return
			}
		}

		var advertisements []models.Advertisement
		query := db.DB.Preload("User").Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

		if statusStr != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&advertisements).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch advertisements",
			})
			return
		}

		var totalRecords int64
		countQuery := db.DB.Model(&models.Advertisement{})
		if statusStr != "" {
			countQuery = countQuery.Where("status = ?", status)
		}
		if err := countQuery.Count(&totalRecords).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to count total records",
			})
			return
		}
		advertisementWithphoto := make([]map[string]interface{}, len(advertisements))
		for i, advertisement := range advertisements {
			photoLink, _ := utilits.GetPhotoLink(advertisement.ImageID)
			advertisementWithphoto[i] = map[string]interface{}{
				"ID":                 advertisement.ID,
				"Text":               advertisement.Text,
				"Status":             advertisement.Status,
				"CreatedAt":          advertisement.CreatedAt,
				"UserID":             advertisement.UserID,
				"photoLink":          photoLink,
				"DeletedFromChannel": advertisement.DeletedFromChannel,
				"UserName":           advertisement.User.Username,
				"FL":                 advertisement.User.FirstName + " " + advertisement.User.LastName,
			}
		}

		totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]interface{}{
				"advertisements": advertisementWithphoto,
				"currentPage":    page,
				"totalPages":     totalPages,
			},
		})
	})
}
