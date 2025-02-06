package main

import "net/http"

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if apiCfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	apiCfg.fileserverHits.Store(0)
	apiCfg.db.ResetUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))

}
