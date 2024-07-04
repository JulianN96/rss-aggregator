package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Initializing Server...")
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	v1Router := chi.NewRouter()

	router.Mount("/v1", v1Router)
	v1Router.Get("/healthcheck", handlerReadiness)
	v1Router.Get("/err", hanlderErr)

	fmt.Printf("Server Starting and listening on port %v...\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
