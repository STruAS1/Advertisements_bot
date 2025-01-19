package routes

import (
	"tgbotBARAHOLKA/backend/routes/Advertisement"
	globalsetings "tgbotBARAHOLKA/backend/routes/GlobalSetings"
	users "tgbotBARAHOLKA/backend/routes/Users"
	"tgbotBARAHOLKA/backend/routes/adsSetings"
	"tgbotBARAHOLKA/backend/routes/payments"

	// "SO/routes/user"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	adsSetings.RegisterAdsSetingsRoutes(r)
	Advertisement.RegisterAdsRoutes(r)
	globalsetings.RegisterAdsRoutes(r)
	payments.RegisterAdsRoutes(r)
	users.RegisterAdsRoutes(r)
}
