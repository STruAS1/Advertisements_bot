package backend

import (
	"fmt"
	"net/http"
	"tgbotBARAHOLKA/backend/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func StartBackend() {
	fmt.Println("Запуск бэкенда")
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
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
