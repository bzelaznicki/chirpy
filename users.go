package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct { ///User struct - global, as may be reused
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct { ///Parameters from request
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err) // parse error
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

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add user", err)
		return
	}

	user := User{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)

}
