package main

import (
	"database/sql"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
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

func (cfg *apiConfig) checkIfUserExistsByUUID(r *http.Request, id uuid.UUID) (bool, error) {
	_, err := cfg.db.GetUserByUUID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
