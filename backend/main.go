package main

import (
	"log"
	"net/http"

	"github.com/deepak-2605/collab-kanban/backend/internal/config"
	"github.com/deepak-2605/collab-kanban/backend/internal/db"
	"github.com/deepak-2605/collab-kanban/backend/internal/handler"
	"github.com/deepak-2605/collab-kanban/backend/internal/middleware"
	"github.com/deepak-2605/collab-kanban/backend/internal/repository"
	"github.com/deepak-2605/collab-kanban/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}

	cfg := config.Load()
	database := db.Connect(cfg.MongoURI, cfg.DBName)

	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/me", authHandler.Me)
	})

	log.Println("Listening on " + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))

}
