package main

import (
	"database/sql"
	"net/http"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {

	chirpId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing Chirp ID", err)
	}

	bearerToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token authentication failed", err)
		return
	}

	userId, err := cfg.validateAccessToken(bearerToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token authentication failed", err)
		return
	}

	chirpToDelete, err := cfg.db.GetSingleChirpByUUID(r.Context(), chirpId)

	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error getting Chirp", err)
		return
	}

	if chirpToDelete.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Access denied", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpToDelete.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting Chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
