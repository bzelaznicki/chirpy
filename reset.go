package main

import "net/http"

func (apiCfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Swap(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("Reset stats successfully"))
	})
}
