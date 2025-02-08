package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
	"github.com/bzelaznicki/chirpy/internal/database"
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

func (cfg *apiConfig) generateUserToken(userId uuid.UUID, timeToExpire time.Duration) (string, error) {
	token, err := auth.MakeJWT(userId, cfg.secret, timeToExpire)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (cfg *apiConfig) generateRefreshToken(r *http.Request, userId uuid.UUID) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		return "", fmt.Errorf("error generating refresh token: %s", err)
	}
	token, err := cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userId,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiration),
	})

	if err != nil {
		return "", fmt.Errorf("error adding refresh token: %s", err)
	}

	return token.Token, nil
}
