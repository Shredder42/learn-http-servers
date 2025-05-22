package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/learn-http-servers/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

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
	db             *database.Queries
	fileserverHits atomic.Int32
}
