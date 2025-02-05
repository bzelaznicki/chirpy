package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		response := fmt.Sprintf("Hits: %d\n", apiCfg.fileserverHits.Load())
		w.Write([]byte(response))
	})
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fileSrv := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.HandleFunc("/healthz", handlerReady)
	mux.Handle("/metrics", apiCfg.handlerMetrics())
	mux.Handle("/reset", apiCfg.handlerReset())

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileSrv))

	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
