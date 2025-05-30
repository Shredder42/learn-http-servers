package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Shredder42/learn-http-servers/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		secret:         secret,
		platform:       platform,
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUsers)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerRetrieveChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirps)
	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serverMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serverMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	serverMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUsers)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirps)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	log.Printf("Serving on port: %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

type apiConfig struct {
	secret         string
	platform       string
	db             *database.Queries
	fileserverHits atomic.Int32
}
