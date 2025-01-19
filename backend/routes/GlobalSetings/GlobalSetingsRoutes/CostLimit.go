package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type CostLimitRequest struct {
	CostLimit uint `json:"CostLimit"`
}

func CostLimit(r chi.Router) {
	r.Get("/CostLimit", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]uint{
				"CostLimit": config.GlobalSettings.Ads.CostLimit,
			},
		})
	})
	r.Put("/CostLimit", func(w http.ResponseWriter, r *http.Request) {
		var CostLimit CostLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&CostLimit); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings

		setings.Ads.CostLimit = CostLimit.CostLimit
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}
