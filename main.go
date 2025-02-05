package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()
	srv := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fileSrv := http.FileServer(http.Dir(filepathRoot))

	mux.Handle("/", fileSrv)
	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
