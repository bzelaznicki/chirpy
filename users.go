package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/bzelaznicki/chirpy/internal/database"
)

type User struct { ///User struct - global, as may be reused
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct { ///Parameters from request
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err) // parse error
		return
	}

	if len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "Enter a password!", fmt.Errorf("no password provided for %s", params.Email))
		return
	}

	if !validateEmail(params.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email address", fmt.Errorf("invalid email address: %s", params.Email))
		return
	}

	exists, err := cfg.checkIfUserExists(r, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unexpected error", err)
		return
	}
	if exists {
		respondWithError(w, http.StatusConflict, "User already exists", fmt.Errorf("User already exists: %s", params.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	createUser := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), createUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
		return
	}
	token, err := cfg.generateUserToken(dbUser.ID, defaultTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}
	user := User{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusCreated, user)

}
