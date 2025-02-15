package main

import (
	"net/http"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpByUUID(r *http.Request, id uuid.UUID) (Chirp, error) {
	chirp, err := cfg.db.GetSingleChirpByUUID(r.Context(), id)
	if err != nil {
		return Chirp{}, err
	}
	convChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	return convChirp, nil
}

func (cfg *apiConfig) validateAccessToken(token string) (uuid.UUID, error) {
	return auth.ValidateJWT(token, cfg.secret)
}
