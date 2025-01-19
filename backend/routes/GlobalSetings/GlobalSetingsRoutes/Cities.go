package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db"
	"tgbotBARAHOLKA/db/models"

	"github.com/go-chi/chi/v5"
)

type AddCityRequest struct {
	Title string `json:"Title"`
}
type CityCountRequest struct {
	Count int8 `json:"Count"`
}

func City(r chi.Router) {
	r.Get("/Cities", func(w http.ResponseWriter, r *http.Request) {
		var Cities []models.Cities
		db.DB.Order("title ASC").Find(&Cities)
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    Cities,
		})
	})
	r.Post("/City", func(w http.ResponseWriter, r *http.Request) {
		var AddCityRequest AddCityRequest
		if err := json.NewDecoder(r.Body).Decode(&AddCityRequest); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		NweCity := models.Cities{
			Title: AddCityRequest.Title,
		}
		db.DB.Create(&NweCity)
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    NweCity,
		})
	})

	r.Delete("/City", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("ID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		if err := db.DB.Where(models.Cities{ID: uint(id)}).Delete(&models.Cities{}).Error; err != nil {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Message: "Record not found",
			})
			return
		}
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "OK",
			Data:    nil,
		})
	})

	r.Get("/CityCount", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]int8{
				"MaxCount": config.GlobalSettings.City.MaxCountOfCity,
			},
		})
	})

	r.Put("/CityCount", func(w http.ResponseWriter, r *http.Request) {
		var Count CityCountRequest
		if err := json.NewDecoder(r.Body).Decode(&Count); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.City.MaxCountOfCity = Count.Count
		config.Save(setings)
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]int8{
				"MaxCount": config.GlobalSettings.City.MaxCountOfCity,
			},
		})
	})
}
