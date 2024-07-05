package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JulianN96/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Initializing Server...")
	//Load env variables
	godotenv.Load()
	//Get port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}
	//Get DB Url 
	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL is not found in the enviroment")
	}
	//Establish connection to DB
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connecting to db", err)
	}

	//Assign query system to struct we can reuse.
	apiCfg := apiConfig{
		DB: database.New(dbConn),
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

	//ROUTES
	router.Mount("/v1", v1Router)
	v1Router.Get("/healthcheck", handlerReadiness)
	v1Router.Get("/err", hanlderErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feedfollows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feedfollows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feedfollows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))


	fmt.Printf("Server Starting and listening on port %v...\n", port)
	bootErr := server.ListenAndServe()
	if bootErr != nil {
		log.Fatal(err)
	}
}
