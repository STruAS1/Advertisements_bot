package backend

import (
	"net/http"
	"tgbotBARAHOLKA/backend/routes"

	"github.com/go-chi/chi/v5"
)

func StartBackend() {
	r := chi.NewRouter()
	routes.RegisterRoutes(r)
	http.ListenAndServe(":8080", r)
}
