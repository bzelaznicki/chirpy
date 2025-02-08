package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.authenticateUser(r, params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}

	token, err := cfg.generateUserToken(dbUser.ID, defaultTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	refreshToken, err := cfg.generateRefreshToken(r, dbUser.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token", err)
		return
	}

	authUser := User{
		ID:           dbUser.ID.String(),
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, authUser)

}
