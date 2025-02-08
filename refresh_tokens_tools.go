package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bzelaznicki/chirpy/internal/auth"
)

func (cfg *apiConfig) validateRefreshToken(r *http.Request) (*RefreshToken, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return nil, err
	}

	// Get the refresh token from the database
	dbToken, err := cfg.db.LookUpRefreshToken(r.Context(), token)
	if err != nil {
		return nil, err
	}

	// Check if token is expired
	if time.Now().After(dbToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// Check if token is revoked
	if dbToken.RevokedAt.Valid {
		return nil, fmt.Errorf("token revoked")
	}

	return &RefreshToken{
		Token:     dbToken.Token,
		UserID:    dbToken.UserID,
		ExpiresAt: dbToken.ExpiresAt,
	}, nil
}
