package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleRedUpgrade(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.EnableChirpyRed(r.Context(), params.Data.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
