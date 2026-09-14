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
	"github.com/go-chi/cors"
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

	boardRepo := repository.NewBoardRepository(database)
	boardService := service.NewBoardService(boardRepo)
	boardHandler := handler.NewBoardHandler(boardService)

	columnRepo := repository.NewColumnRepository(database)
	columnService := service.NewColumnService(columnRepo, boardRepo)
	columnHandler := handler.NewColumnHandler(columnService)

	cardRepo := repository.NewCardRepository(database)
	cardService := service.NewCardService(cardRepo, columnRepo, boardRepo)
	cardHandler := handler.NewCardHandler(cardService)

	r := chi.NewRouter()

	// allow the React dev server (different origin) to call this API
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/me", authHandler.Me)

		r.Get("/boards", boardHandler.List)
		r.Post("/boards", boardHandler.Create)
		r.Patch("/boards/{boardID}", boardHandler.Rename)
		r.Delete("/boards/{boardID}", boardHandler.Delete)
		r.Get("/boards/{boardID}/columns", columnHandler.List)
		r.Post("/boards/{boardID}/columns", columnHandler.Create)
		r.Patch("/columns/{columnID}", columnHandler.Rename)
		r.Delete("/columns/{columnID}", columnHandler.Delete)

		r.Get("/columns/{columnID}/cards", cardHandler.List)
		r.Post("/columns/{columnID}/cards", cardHandler.Create)
		r.Patch("/cards/{cardID}", cardHandler.Update)
		r.Patch("/cards/{cardID}/move", cardHandler.Move)
		r.Delete("/cards/{cardID}", cardHandler.Delete)
	})

	log.Println("Listening on " + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))

}
