package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))

}

var usersColl *mongo.Collection

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env vars")
	}
	r := chi.NewRouter()

	client := connectMongo()
	usersColl = client.Database("kanban").Collection("users")

	r.Get("/health", healthHandler)
	r.Post("/register", registerHandler)
	r.Post("/login", loginHandler)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware) // everything in here requires a valid token
		r.Get("/me", meHandler)
	})
	port := os.Getenv("PORT")
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
