package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/bzelaznicki/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Missing or invalid bearer token", err)
		return
	}

	userId, err := cfg.authenticateUserByToken(bearerToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token authentication failed", err)
		return
	}

	valid, msg := validateLength(params.Body)

	if !valid {
		respondWithError(w, http.StatusBadRequest, msg, nil)
		return
	}

	cleanedChirp, err := profanityChecker(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to proces chirp.", err)
		return
	}

	chirp, err := cfg.db.PostChirp(r.Context(), database.PostChirpParams{
		Body:   cleanedChirp,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to post Chirp", err)
		return
	}

	postedChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, postedChirp)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting Chirps", err)
		return
	}

	chirpsResponse := []Chirp{}

	for _, chirp := range chirps {
		transformedChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpsResponse = append(chirpsResponse, transformedChirp)
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {

	chirpId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing ID", err)
		return
	}

	chirp, err := cfg.getChirpByUUID(r, chirpId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error getting Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)

}
