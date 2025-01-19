package adsSetings

import (
	"tgbotBARAHOLKA/backend/routes/adsSetings/adsSetingsRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsSetingsRoutes(r chi.Router) {
	r.Route("/adsSetings", func(r chi.Router) {
		adsSetingsRoutes.Input(r)
		adsSetingsRoutes.Type(r)
		adsSetingsRoutes.Sufix(r)
	})
}
