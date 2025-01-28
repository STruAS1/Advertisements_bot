package verification

import (
	verifictionroutes "tgbotBARAHOLKA/backend/routes/Verification/VerifictionRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsRoutes(r chi.Router) {
	r.Route("/Verifiaction", func(r chi.Router) {
		verifictionroutes.Aplications(r)
	})
}
