package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedUUID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUUID)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "secret"
	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	assert.Error(t, err)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := -time.Hour // Token already expired

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ValidateJWT(token, tokenSecret)
	assert.Error(t, err)

}

func TestValidateJWT_WrongSigningMethod(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret"
	expiresIn := time.Hour

	// Create a token with a different signing method
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	token, err := newToken.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = ValidateJWT(token, tokenSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

func TestValidateJWT_WrongSecretKey(t *testing.T) {
	userId := uuid.New()
	originalSecret := "secret"
	wrongSecret := "wrong"
	expiresIn := time.Hour

	newToken, err := MakeJWT(userId, originalSecret, expiresIn)
	assert.NoError(t, err)
	_, err = ValidateJWT(newToken, wrongSecret)

	assert.Error(t, err)

}
