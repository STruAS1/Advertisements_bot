package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type DomainRequest struct {
	Domain string `json:"Domain"`
}

func AllowedDomines(r chi.Router) {
	r.Get("/Domines", func(w http.ResponseWriter, r *http.Request) {
		DominesWithIndex := make([]map[string]interface{}, len(config.GlobalSettings.WitheListDomines))
		for i, Domain := range config.GlobalSettings.WitheListDomines {
			DominesWithIndex[i] = map[string]interface{}{
				"index":  i,
				"Domain": Domain,
			}
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    DominesWithIndex,
		})
	})
	r.Post("/Domain", func(w http.ResponseWriter, r *http.Request) {
		var Domain DomainRequest
		if err := json.NewDecoder(r.Body).Decode(&Domain); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.WitheListDomines = append(setings.WitheListDomines, Domain.Domain)
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Delete("/Domain", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("Index")
		Index, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		if int(Index) >= len(setings.WitheListDomines) {
			http.Error(w, "Index out of range", http.StatusBadRequest)
			return
		}
		setings.WitheListDomines = append(setings.WitheListDomines[:Index], setings.WitheListDomines[Index+1:]...)
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}
