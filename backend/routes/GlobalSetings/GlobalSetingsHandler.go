package globalsetings

import (
	globalsetingsroutes "tgbotBARAHOLKA/backend/routes/GlobalSetings/GlobalSetingsRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsRoutes(r chi.Router) {
	r.Route("/setings", func(r chi.Router) {
		globalsetingsroutes.MainText(r)
		globalsetingsroutes.City(r)
		globalsetingsroutes.Payments(r)
		globalsetingsroutes.AdsText(r)
		globalsetingsroutes.Docs(r)
		globalsetingsroutes.CostLimit(r)
		globalsetingsroutes.Buttons(r)
	})
}
