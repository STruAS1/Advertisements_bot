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

type UpdateStatusRequest struct {
	Status uint8 `json:"Status"`
}

func UpdateStatus(r chi.Router) {
	r.Put("/UpdateStatus", func(w http.ResponseWriter, r *http.Request) {
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
		var Payment models.Payments
		db.DB.Preload("User").Where(models.Payments{ID: uint(ID)}).First(&Payment)
		if Payment.Status != 0 {
			http.Error(w, "Failed to update record", http.StatusInternalServerError)
			return
		}
		if UpdateStatus.Status == 1 {
			if err := db.DB.Model(&models.Payments{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status": UpdateStatus.Status,
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			if err := db.DB.Model(&models.User{}).
				Where("id = ?", uint(Payment.UserID)).
				Updates(map[string]interface{}{
					"balance": uint(Payment.User.Balance + Payment.Amount),
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			text := "💸Ваша заявка на пополнение " + strconv.Itoa(int(Payment.Amount)) + "₩ принята!"
			utilits.SendMessageToUser(text, Payment.User.TelegramID)
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		} else if UpdateStatus.Status == 2 {
			if err := db.DB.Model(&models.Payments{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status": UpdateStatus.Status,
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			text := "🚫Ваша заявка на пополнение " + strconv.Itoa(int(Payment.Amount)) + "₩ отклонена!"
			utilits.SendMessageToUser(text, Payment.User.TelegramID)
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		}
	})
}
