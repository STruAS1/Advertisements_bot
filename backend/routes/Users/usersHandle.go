package users

import (
	usersroutes "tgbotBARAHOLKA/backend/routes/Users/UsersRoutes"

	"github.com/go-chi/chi/v5"
)

func RegisterAdsRoutes(r chi.Router) {
	r.Route("/Users", func(r chi.Router) {
		usersroutes.GetAllUsers(r)
		usersroutes.BanRoutes(r)
	})
}
