package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/bzelaznicki/chirpy/internal/database"
)

func (cfg *apiConfig) authenticateUser(r *http.Request, email, password string) (database.User, error) {
	user, err := cfg.db.GetUser(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return database.User{}, fmt.Errorf("invalid credentials")
		}
		return database.User{}, err
	}
	err = auth.CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return database.User{}, fmt.Errorf("invalid credentials")
	}
	return user, nil

}
