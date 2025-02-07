package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.ExpiresInSeconds < 1 || params.ExpiresInSeconds > defaultTokenExpiration {
		params.ExpiresInSeconds = defaultTokenExpiration
	}

	dbUser, err := cfg.authenticateUser(r, params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}

	token, err := cfg.generateUserToken(dbUser.ID, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	authUser := User{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, authUser)

}
