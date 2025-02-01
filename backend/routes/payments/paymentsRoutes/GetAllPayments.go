package paymentsroutes

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

func GetAllPayments(r chi.Router) {
	r.Get("/payments", func(w http.ResponseWriter, r *http.Request) {

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

		var payments []models.Payments
		query := db.DB.Preload("User").Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

		if statusStr != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&payments).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch advertisements",
			})
			return
		}

		var totalRecords int64
		countQuery := db.DB.Model(&models.Payments{})
		if statusStr != "" {
			countQuery = countQuery.Where("status = ?", status)
		}
		if err := countQuery.Count(&totalRecords).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to count total records",
			})
			return
		}
		paymentSorted := make([]map[string]interface{}, len(payments))
		for i, payment := range payments {
			photoLink, _ := utilits.GetPhotoLink(payment.PhotoUrl)
			paymentSorted[i] = map[string]interface{}{
				"ID":        payment.ID,
				"Metod":     payment.Metod,
				"Status":    payment.Status,
				"CreatedAt": payment.CreatedAt,
				"amount":    payment.Amount,
				"UserID":    payment.UserID,
				"UserName":  payment.User.Username,
				"FL":        payment.User.FirstName + " " + payment.User.LastName,
				"photoLink": photoLink,
			}
		}

		totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]interface{}{
				"advertisements": paymentSorted,
				"currentPage":    page,
				"totalPages":     totalPages,
			},
		})
	})
	r.Delete("/payment", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		db.DB.Delete(&models.Payments{}, ID)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}
