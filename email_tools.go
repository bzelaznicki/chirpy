package main

import (
	"database/sql"
	"net/http"
	"net/mail"
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil

}

func (cfg *apiConfig) checkIfUserExists(r *http.Request, email string) (bool, error) {
	_, err := cfg.db.GetUser(r.Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // User doesn't exist
		}
		return false, err // Some other error occurred
	}
	return true, nil // User exists
}
