package payments

import (
	paymentsroutes "tgbotBARAHOLKA/backend/routes/payments/paymentsRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		paymentsroutes.GetAllPayments(r)
		paymentsroutes.UpdateStatus(r)
	})
}
