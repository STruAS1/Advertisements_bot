package adsSetings

import (
	"tgbotBARAHOLKA/backend/routes/adsSetings/adsSetingsRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsSetingsRoutes(r chi.Router) {
	r.Route("/adsSetings", func(r chi.Router) {
		adsSetingsRoutes.GetAllAdvertisementTypes(r)
		adsSetingsRoutes.GetAllAdvertisementInputs(r)
		adsSetingsRoutes.GetAdvertisementInputInfo(r)
		adsSetingsRoutes.UpdateAdvertisementInput(r)
		adsSetingsRoutes.UpdateAdvertisementInput(r)
		adsSetingsRoutes.UpdateAdvertisementType(r)
		adsSetingsRoutes.GetAdvertisementTypeInfo(r)
		adsSetingsRoutes.CreateAdvertisementInput(r)
		adsSetingsRoutes.CreateAdvertisementType(r)
	})
}
