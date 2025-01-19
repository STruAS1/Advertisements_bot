package backend

import (
	"net/http"
	"tgbotBARAHOLKA/backend/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func StartBackend() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	routes.RegisterRoutes(r)
	http.ListenAndServe(":8080", r)
}
