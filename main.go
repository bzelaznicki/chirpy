package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fileSrv := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.HandleFunc("/healthz", handlerReady)

	mux.Handle("/app/", fileSrv)

	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
