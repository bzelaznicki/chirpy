package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/bzelaznicki/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid bearer token", err)
		return
	}

	userId, err := cfg.validateAccessToken(bearerToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token authentication failed", err)
		return
	}

	err = validateUserDetails(params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid credentials, check email and password", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	updateParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now().UTC(),
		ID:             userId,
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), updateParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user", err)
		return
	}

	dbUser := response{
		ID:        updatedUser.ID.String(),
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, dbUser)

}
