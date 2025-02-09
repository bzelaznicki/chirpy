package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/bzelaznicki/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		response := fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", apiCfg.fileserverHits.Load())
		w.Write([]byte(response))
	})
}

func main() {
	godotenv.Load()
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")

	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")

	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalf("database connection failed %s", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
	}

	fileSrv := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetSingleChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUserUpdate)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handleChirpDelete)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileSrv))

	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
