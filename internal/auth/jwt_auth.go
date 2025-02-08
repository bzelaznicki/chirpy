package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	issueDate := jwt.NewNumericDate(time.Now().UTC())
	expirationDate := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))

	newToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  issueDate,
			ExpiresAt: expirationDate,
			Subject:   userID.String(),
		},
	)
	signedString, err := newToken.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %s", err)
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
			}
			return []byte(tokenSecret), nil
		},
		// Add explicit validation options
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error processing token: %v", err)
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	// Explicit expiration check
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return uuid.UUID{}, fmt.Errorf("token has expired")
	}

	subject := claims.Subject
	parsedUUID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing UUID: %s", err)
	}

	return parsedUUID, nil
}
