package Advertisement

import (
	adsroutes "tgbotBARAHOLKA/backend/routes/Advertisement/AdsRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsRoutes(r chi.Router) {
	r.Route("/ads", func(r chi.Router) {
		adsroutes.GetAllAdvertisements(r)
	})
}
