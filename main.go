package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"mondragon-ai/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main () {
	fmt.Println("starting go server....")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the database URL
	dbURL := os.Getenv("DB_URL")

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create SQLC database queries
	dbQueries := database.New(db)

	fileServer := http.FileServer(http.Dir("."))

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: "dev",
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/", apiCfg.handlerGetChirpByID)


	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
