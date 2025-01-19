package globalsetingsroutes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tgbotBARAHOLKA/config"

	"github.com/go-chi/chi/v5"
)

type PaymentCrateRequest struct {
	Title       string `json:"Title"`
	Cardnumber  string `json:"Cardnumber"`
	Description string `json:"Description"`
}

type LimiOfAmountRequest struct {
	MaxAmount uint `json:"Max"`
	MinAmount uint `json:"Min"`
}

func Payments(r chi.Router) {
	r.Get("/Payments", func(w http.ResponseWriter, r *http.Request) {
		paymentsWithIndex := make([]map[string]interface{}, len(config.GlobalSettings.Payments.Metods))
		for i, method := range config.GlobalSettings.Payments.Metods {
			paymentsWithIndex[i] = map[string]interface{}{
				"index":       i,
				"Title":       method.Title,
				"Cardnumber":  method.Title,
				"Description": method.Discription,
			}
		}

		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data:    paymentsWithIndex,
		})
	})
	r.Post("/Payment", func(w http.ResponseWriter, r *http.Request) {
		var Payment PaymentCrateRequest
		if err := json.NewDecoder(r.Body).Decode(&Payment); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Payments.Metods = append(setings.Payments.Metods, config.PaymentsMetod{Title: Payment.Title, Discription: Payment.Description, Cardnumber: Payment.Cardnumber})
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})

	})
	r.Delete("/Payment", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		idStr := queryParams.Get("Index")
		Index, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid ID: must be a positive integer", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		if int(Index) >= len(setings.Payments.Metods) {
			http.Error(w, "Index out of range", http.StatusBadRequest)
			return
		}
		setings.Payments.Metods = append(setings.Payments.Metods[:Index], setings.Payments.Metods[Index+1:]...)
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
	r.Get("/LimiOfAmount", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, SuccessResponse{
			Message: "ok",
			Data: map[string]uint{
				"Max": config.GlobalSettings.Payments.MaxAmount,
				"Min": config.GlobalSettings.Payments.MinimalAmount,
			},
		})
	})
	r.Put("/LimiOfAmount", func(w http.ResponseWriter, r *http.Request) {
		var LimiOfAmount LimiOfAmountRequest
		if err := json.NewDecoder(r.Body).Decode(&LimiOfAmount); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		setings := config.GlobalSettings
		setings.Payments.MinimalAmount = LimiOfAmount.MinAmount
		setings.Payments.MaxAmount = LimiOfAmount.MaxAmount
		config.Save(setings)
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "Ok",
		})
	})
}
