package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("error generating random key: %s", err)
	}
	randomString := hex.EncodeToString(bytes)

	return randomString, nil
}
