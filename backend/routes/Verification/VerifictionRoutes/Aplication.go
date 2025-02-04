package verifictionroutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"
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

type UpdateStatusRequest struct {
	Status uint8  `json:"Status"`
	Msg    string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Aplications(r chi.Router) {
	r.Get("/Aplications", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		limitStr := queryParams.Get("limit")
		pageStr := queryParams.Get("page")
		statusStr := queryParams.Get("status")

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

		var Aplications []models.VerificationAplication
		query := db.DB.Preload("User").Order("created_at desc").Limit(limit).Offset((page - 1) * limit)
		if statusStr != "" {
			query = query.Where("status = ?", status)
		}
		if err := query.Find(&Aplications).Error; err != nil {
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to fetch advertisements",
			})
			return
		}

		var totalRecords int64
		countQuery := db.DB.Model(&models.VerificationAplication{})

		if err := countQuery.Count(&totalRecords).Error; err != nil {
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to count total records",
			})
			return
		}
		AplicationsWithphoto := make([]map[string]interface{}, len(Aplications))
		for i, Aplication := range Aplications {
			CardIdFileUrl, _ := utilits.GetPhotoLink(Aplication.CardIdFileID)
			DocumentFileUrl, _ := utilits.GetPhotoLink(Aplication.DocumentFileID)
			AplicationsWithphoto[i] = map[string]interface{}{
				"ID": Aplication.ID,
				"AplicationFirstNameAndAplicationLastNameAndAplicationPatronymic": Aplication.Patronymic + " " + Aplication.FirstName + " " + Aplication.LastName,
				"VisaType":        Aplication.VisaType,
				"Services":        Aplication.Services,
				"Status":          Aplication.Status,
				"CreatedAt":       Aplication.CreatedAt,
				"UserID":          Aplication.UserID,
				"CardIdFileUrl":   CardIdFileUrl,
				"DocumentFileUrl": DocumentFileUrl,
				"UserName":        Aplication.User.Username,
				"FL":              Aplication.User.FirstName + " " + Aplication.User.LastName,
				"phone":           Aplication.User.Phone,
			}
		}
		totalPages := (totalRecords + int64(limit) - 1) / int64(limit)

		WriteJSON(w, http.StatusOK, SuccessResponse{
			Message: "Ok",
			Data: map[string]interface{}{
				"Aplications": AplicationsWithphoto,
				"currentPage": page,
				"totalPages":  totalPages,
				"totalUsers":  totalRecords,
			},
		})
	})
	r.Delete("/Aplication", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var Aplication models.VerificationAplication
		db.DB.Where(&models.VerificationAplication{ID: uint(ID)}).Find(&Aplication)
		if err := db.DB.Model(&models.User{}).
			Where("id = ?", uint(Aplication.UserID)).
			Updates(map[string]interface{}{
				"verification": false,
			}).Error; err != nil {
			WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
				Message: "Failed to update record",
			})
			return
		}
		db.DB.Delete(&models.VerificationAplication{}, ID)
		WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
	r.Put("/Aplication", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		ID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		var UpdateStatus UpdateStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&UpdateStatus); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		var Aplication models.VerificationAplication
		db.DB.Preload("User").Where(models.VerificationAplication{ID: uint(ID)}).First(&Aplication)
		if Aplication.Status != 0 {
			http.Error(w, "Failed to update record", http.StatusInternalServerError)
			return
		}
		if UpdateStatus.Status == 1 {
			if err := db.DB.Model(&models.VerificationAplication{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status": UpdateStatus.Status,
				}).Error; err != nil {
				WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
					Message: "Failed to update record",
				})
				return
			}
			if err := db.DB.Model(&models.User{}).
				Where("id = ?", uint(Aplication.UserID)).
				Updates(map[string]interface{}{
					"verification": true,
					"balance":      Aplication.User.Balance - config.GlobalSettings.VerificationCost,
				}).Error; err != nil {
				WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
					Message: "Failed to update record",
				})
				return
			}
			text := "✅ Ваша заявка на верификацию успешно подтверждена!"
			if UpdateStatus.Msg != "" {
				text += fmt.Sprintf("\n\n%s", UpdateStatus.Msg)
			}
			utilits.SendMessageToUser(text, int64(Aplication.User.TelegramID))
			WriteJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		} else if UpdateStatus.Status == 2 {
			if err := db.DB.Model(&models.VerificationAplication{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status": UpdateStatus.Status,
				}).Error; err != nil {
				WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
					Message: "Failed to update record",
				})
				return
			}
			text := "❌ К сожалению, ваша заявка на верификацию была отклонена."
			if UpdateStatus.Msg != "" {
				text += fmt.Sprintf("\n\n%s", UpdateStatus.Msg)
			}
			utilits.SendMessageToUser(text, int64(Aplication.User.TelegramID))
			WriteJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		}
	})
}
