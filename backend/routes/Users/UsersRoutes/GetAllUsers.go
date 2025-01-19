package usersroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func GetAllUsers(r chi.Router) {
	r.Get("/Users", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		limitStr := queryParams.Get("limit")
		pageStr := queryParams.Get("page")

		var limit, page int
		var err error

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

		var users []models.User
		query := db.DB.Order("created_at desc").Limit(limit).Offset((page - 1) * limit)

		if err := query.Find(&users).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch advertisements",
			})
			return
		}

		var totalRecords int64
		countQuery := db.DB.Model(&models.User{})

		if err := countQuery.Count(&totalRecords).Error; err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to count total records",
			})
			return
		}
		UsersWithCount := make([]map[string]interface{}, len(users))
		for i, user := range users {
			var CountOfAds int64
			db.DB.Model(&models.Advertisement{}).Where(&models.Advertisement{UserID: user.ID}).Count(&CountOfAds)
			UsersWithCount[i] = map[string]interface{}{
				"ID":         user.ID,
				"Balance":    user.Balance,
				"CountOfAds": CountOfAds,
				"CreatedAt":  user.CreatedAt,
				"UserName":   user.Username,
				"FL":         user.FirstName + " " + user.LastName,
			}
		}

		totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]interface{}{
				"advertisements": UsersWithCount,
				"currentPage":    page,
				"totalPages":     totalPages,
				"totalUsers":     totalRecords,
			},
		})
	})
}
