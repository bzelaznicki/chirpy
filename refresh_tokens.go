package main

import (
	"net/http"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := cfg.validateRefreshToken(r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	authToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, defaultTokenExpiration)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	resp := response{
		Token: authToken,
	}

	respondWithJSON(w, http.StatusOK, resp)

}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := cfg.validateRefreshToken(r)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
