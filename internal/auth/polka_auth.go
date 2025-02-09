package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	header := strings.Split(headers.Get("Authorization"), " ")

	if len(header) < 2 {
		return "", fmt.Errorf("no valid auth header")
	}

	if header[0] != "ApiKey" {
		return "", fmt.Errorf("no valid API Key")
	}

	apiKey := strings.TrimSpace(header[1])

	if len(apiKey) == 0 {
		return "", fmt.Errorf("API Key missing")
	}

	return apiKey, nil

}
