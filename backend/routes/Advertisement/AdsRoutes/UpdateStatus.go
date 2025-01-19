package adsroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"
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
		var AD models.Advertisement
		db.DB.Preload("User").Where(models.Advertisement{ID: uint(ID)}).First(&AD)
		if AD.Status != 0 {
			http.Error(w, "Failed to update record", http.StatusInternalServerError)
			return
		}
		if UpdateStatus.Status == 1 {
			var msgText string = AD.Text
			msgText += "\n👉<b><a href='https://t.me/" + AD.User.Username + "'>Написать автору</a></b>👈"
			msgText += "\n\n" + config.GlobalSettings.Ads.Sufix
			msgId := utilits.SendMessageToChnale(msgText, AD.ImageID)
			if err := db.DB.Model(&models.Advertisement{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status":    UpdateStatus.Status,
					"massge_id": msgId,
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		} else if UpdateStatus.Status == 2 {
			if err := db.DB.Model(&models.Advertisement{}).
				Where("id = ?", uint(ID)).
				Updates(map[string]interface{}{
					"status": UpdateStatus.Status,
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			if err := db.DB.Model(&models.User{}).
				Where("id = ?", uint(AD.User.ID)).
				Updates(map[string]interface{}{
					"balance": uint(AD.User.Balance + AD.CostUser),
				}).Error; err != nil {
				http.Error(w, "Failed to update record", http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "Ok",
			})
		}
	})
}
