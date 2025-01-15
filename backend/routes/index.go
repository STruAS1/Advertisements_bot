package routes

import (
	"tgbotBARAHOLKA/backend/routes/Advertisement"
	"tgbotBARAHOLKA/backend/routes/adsSetings"

	// "SO/routes/user"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	adsSetings.RegisterAdsSetingsRoutes(r)
	Advertisement.RegisterAdsRoutes(r)
}
